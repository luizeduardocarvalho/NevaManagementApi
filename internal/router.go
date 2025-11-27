package internal

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/luizeduardocarvalho/labflux-functions/functions"
	custommiddleware "github.com/luizeduardocarvalho/labflux-functions/pkg/middleware"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"status": "ok"})
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public routes (no auth required)
		r.Group(func(r chi.Router) {
			// Health checks
			r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
				render.JSON(w, r, map[string]string{"message": "pong"})
			})

			// Authentication endpoints
			r.Post("/auth/signup", functions.SignUp)
			r.Post("/auth/signin", functions.SignIn)
			r.Post("/auth/login", functions.SignIn) // Alias for signin

			// Clerk webhook
			r.Post("/webhooks/clerk", functions.ClerkWebhookHandler)

			// Invitation acceptance (requires Clerk user ID)
			r.Post("/auth/invite/{token}", functions.AcceptInvitation)

			// Public invitation endpoints
			r.Get("/invitations/token/{token}", functions.GetInvitationByToken)
		})

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.AuthMiddleware)

			// User management
			r.Get("/auth/me", functions.GetCurrentUser)
			r.Post("/auth/refresh", functions.RefreshToken)

			// Organization management
			functions.RegisterOrganizationRoutes(r)

			// Laboratory management
			r.Post("/laboratories", functions.CreateLaboratory)
			r.Get("/laboratories/{laboratoryId}", functions.GetLaboratory)

			// Invitation management
			r.Post("/laboratories/{laboratoryId}/invitations", functions.CreateInvitation)
			r.Get("/laboratories/{laboratoryId}/invitations", functions.GetInvitations)
			r.Post("/invitations/{id}/resend", functions.ResendInvitation)
			r.Post("/invitations/{id}/cancel", functions.CancelInvitation)
			r.Delete("/invitations/{id}", functions.DeleteInvitation)

			// Location management
			functions.RegisterLocationRoutes(r)

			// Product management
			functions.RegisterProductRoutes(r)

			// Equipment management (includes usage routes)
			functions.RegisterEquipmentRoutes(r)

			// Researcher management
			functions.RegisterResearcherRoutes(r)

			// Sample management
			functions.RegisterSampleRoutes(r)

			// Replica management
			functions.RegisterReplicaRoutes(r)

			// Routine management
			functions.RegisterRoutineRoutes(r)
		})
	})

	return r
}
