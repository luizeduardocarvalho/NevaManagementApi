// Package labflux contains Cloud Function entry points
package labflux

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/luizeduardocarvalho/labflux-functions/docs"
	"github.com/luizeduardocarvalho/labflux-functions/internal"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
)

func init() {
	if err := database.Initialize(); err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	functions.HTTP("HandleRequest", HandleRequest)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	router := internal.SetupRouter()
	router.ServeHTTP(w, r)
}
