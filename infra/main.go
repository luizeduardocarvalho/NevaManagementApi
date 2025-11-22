package main

import (
	"fmt"

	cloudfunctionsv2 "github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/cloudfunctionsv2"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		cfgGcp := config.New(ctx, "gcp")

		// Config (same keys you already used)
		projectID := cfgGcp.Require("project")
		region := cfgGcp.Get("region")
		if region == "" {
			region = "us-central1"
		}

		databaseURL := cfg.RequireSecret("database-url")
		clerkPublishableKey := cfg.Require("clerk-publishable-key")
		clerkSecretKey := cfg.RequireSecret("clerk-secret-key")
		clerkWebhookSecret := cfg.RequireSecret("clerk-webhook-secret")

		// Bucket for source
		bucket, err := storage.NewBucket(ctx, "labflux-functions-bucket", &storage.BucketArgs{
			Location:                 pulumi.String(region),
			Name:                     pulumi.Sprintf("%s-labflux-functions", projectID),
			UniformBucketLevelAccess: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create archive with only necessary files (exclude infra/ to avoid go.mod conflicts)
		archive := pulumi.NewAssetArchive(map[string]interface{}{
			"function.go":   pulumi.NewFileAsset("../function.go"),
			"go.mod":        pulumi.NewFileAsset("../go.mod"),
			"go.sum":        pulumi.NewFileAsset("../go.sum"),
			"functions":     pulumi.NewFileArchive("../functions"),
			"docs":          pulumi.NewFileArchive("../docs"),
			"internal":      pulumi.NewFileArchive("../internal"),
			"pkg":           pulumi.NewFileArchive("../pkg"),
			".gcloudignore": pulumi.NewFileAsset("../.gcloudignore"),
		})

		srcObj, err := storage.NewBucketObject(ctx, "labflux-functions-source", &storage.BucketObjectArgs{
			Bucket: bucket.Name,
			Name:   pulumi.String("labflux-functions.zip"),
			Source: archive,
		})
		if err != nil {
			return err
		}

		// Gen 2 Cloud Function (HTTP)
		fn, err := cloudfunctionsv2.NewFunction(ctx, "labflux-api", &cloudfunctionsv2.FunctionArgs{
			Name:     pulumi.String("labflux-api"),
			Location: pulumi.String(region),
			Labels: pulumi.StringMap{
				"environment": pulumi.String("dev"),
			},

			BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
				Runtime:    pulumi.String("go122"),         // stay on supported Go runtime
				EntryPoint: pulumi.String("HandleRequest"), // your function symbol
				Source: &cloudfunctionsv2.FunctionBuildConfigSourceArgs{
					StorageSource: &cloudfunctionsv2.FunctionBuildConfigSourceStorageSourceArgs{
						Bucket: bucket.Name,
						Object: srcObj.Name,
					},
				},
				EnvironmentVariables: pulumi.StringMap{
					"GOPROXY":     pulumi.String("https://proxy.golang.org"),
					"GO111MODULE": pulumi.String("on"),
				},
			},

			ServiceConfig: &cloudfunctionsv2.FunctionServiceConfigArgs{
				// Memory/timeout, ingress, env vars live here in Gen 2
				AvailableMemory:            pulumi.String("256M"),
				TimeoutSeconds:             pulumi.Int(60),
				IngressSettings:            pulumi.String("ALLOW_ALL"),
				AllTrafficOnLatestRevision: pulumi.Bool(true),
				EnvironmentVariables: pulumi.StringMap{
					"DATABASE_URL":          databaseURL,
					"CLERK_PUBLISHABLE_KEY": pulumi.String(clerkPublishableKey),
					"CLERK_SECRET_KEY":      clerkSecretKey,
					"CLERK_WEBHOOK_SECRET":  clerkWebhookSecret,
				},
			},
		})
		if err != nil {
			return err
		}

		// Make function publicly accessible via Cloud Functions IAM binding
		_, err = cloudfunctionsv2.NewFunctionIamMember(ctx, "labflux-api-invoker", &cloudfunctionsv2.FunctionIamMemberArgs{
			Project:       pulumi.String(projectID),
			Location:      pulumi.String(region),
			CloudFunction: fn.Name,
			Role:          pulumi.String("roles/cloudfunctions.invoker"),
			Member:        pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}

		// Ensure underlying Cloud Run service stays public as well
		serviceName := fn.ServiceConfig.ApplyT(func(sc *cloudfunctionsv2.FunctionServiceConfig) (string, error) {
			if sc == nil || sc.Service == nil {
				return "", fmt.Errorf("function service config not available")
			}
			return *sc.Service, nil
		}).(pulumi.StringOutput)

		_, err = cloudrun.NewIamMember(ctx, "labflux-api-run-invoker", &cloudrun.IamMemberArgs{
			Project:  pulumi.String(projectID),
			Location: pulumi.String(region),
			Service:  serviceName,
			Role:     pulumi.String("roles/run.invoker"),
			Member:   pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}

		// Export function URL (ServiceConfig includes .Uri in Gen 2)
		url := fn.ServiceConfig.ApplyT(func(sc *cloudfunctionsv2.FunctionServiceConfig) string {
			if sc == nil || sc.Uri == nil {
				return ""
			}
			return *sc.Uri
		}).(pulumi.StringOutput)

		ctx.Export("functionUrl", url)
		ctx.Export("bucketName", bucket.Name)

		return nil
	})
}
