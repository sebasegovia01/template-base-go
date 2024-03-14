package core

import (
	"encoding/json"
	"net/http"
	"template-base-go/src/handlers"
	"template-base-go/src/utils"

	_ "template-base-go/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Routes interface {
	Router() http.Handler
}

type api struct {
	router http.Handler
}

func NewApi(handlers *handlers.Container, logger *utils.Logger) Routes {
	a := &api{}
	a.init(handlers, logger)
	return a
}

// return router
func (a *api) Router() http.Handler {
	return a.router
}

// Init defines and set API routes (void)
func (a *api) init(handlers *handlers.Container, logger *utils.Logger) {

	router := mux.NewRouter()

	// set responses to JSON
	router.Use(JSONContentTypeMiddleware)

	// NotFoundHandler JSON response
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Resource not found",
		})
	})

	// MethodNotAllowedHandler JSON response
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
	})

	//Swagger
	router.PathPrefix("/api-docs").Handler(httpSwagger.WrapHandler)

	subrouter := router.PathPrefix("/service-otp/v1").Subrouter()

	subrouter.Use(LogRequest(logger))

	subrouter.Handle("/generate", http.HandlerFunc(handlers.OtpHandler.Generate)).Methods(http.MethodGet)

	a.router = router
}
