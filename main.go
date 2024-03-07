package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"template-base-go/src/core"
	"template-base-go/src/handlers"
	"template-base-go/src/services"
	"template-base-go/src/utils"

	gorillaHandlers "github.com/gorilla/handlers"
)

// utils container init
func newUtilsContainer() *utils.Container {
	env := utils.NewEnvironment()

	return utils.NewContainer(env)
}

// service container init
func newServicesContainer() *services.Container {

	envUtils := &utils.Environment{} // only to use the lookup method which does not require logger
	logger := &utils.Logger{}

	return services.NewContainer(
		logger,
		envUtils,
	)
}

// handler container init
func newHandlersContainer(servicesCont *services.Container) *handlers.Container {

	exampleHandler := handlers.NewExampleHandler(
		servicesCont.Logger,
		&services.ExampleService{},
	)

	return handlers.NewContainer(exampleHandler)
}

func main() {
	// environment arg
	args := os.Args[1:]

	var localEnv string
	if len(args) > 0 {
		localEnv = os.Args[1:][0] // just 1 arg is received

		// only dev, prod and qa allowed
		if localEnv != "dev" && localEnv != "prod" && localEnv != "qa" {
			localEnv = ""
		}
	}

	if err := utils.ParseEnvironmentFile(localEnv); err != nil {
		log.Fatal(err)
	}

	// CORS
	headers := gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Internal", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With"})
	methods := gorillaHandlers.AllowedMethods([]string{"DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"})
	origins := gorillaHandlers.AllowedOrigins([]string{"*"})

	// init containers
	utilsContainer := newUtilsContainer()

	servicesContainer := newServicesContainer()

	handlersContainers := newHandlersContainer(servicesContainer)

	api := core.NewApi(handlersContainers, &utils.Logger{})

	servicesContainer.Logger.Info("Initializing server")

	port := fmt.Sprintf(":%v", utilsContainer.Environment.GetEnvVar("PORT"))

	if port == ":" {
		port = ":5000"
	}
	servicesContainer.Logger.Info("Server url: http://localhost:3000")

	servicesContainer.Logger.Info("Swagger url: http://localhost:3000/api-docs/index.html")

	log.Fatal(http.ListenAndServe(port, gorillaHandlers.CORS(headers, methods, origins)(api.Router())))
}
