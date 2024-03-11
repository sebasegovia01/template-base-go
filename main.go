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
	utilsContainer := newUtilsContainer()
	servicesContainer := newServicesContainer()

	startHTTPServer(utilsContainer, servicesContainer)

}

func startHTTPServer(utilsContainer *utils.Container, servicesContainer *services.Container) {
	args := os.Args[1:]
	localEnv := ""
	if len(args) > 0 {
		localEnv = args[0]
	}

	if err := utils.ParseEnvironmentFile(localEnv); err != nil {
		log.Fatal(err)
	}

	api := setupAPI()
	logServerInfo(utilsContainer, servicesContainer)

	port := getServerPort(utilsContainer)
	log.Fatal(http.ListenAndServe(port, setupCORS(api.Router())))
}

func setupAPI() core.Routes {
	servicesContainer := newServicesContainer()
	handlersContainers := newHandlersContainer(servicesContainer)
	return core.NewApi(handlersContainers, &utils.Logger{})
}

func setupCORS(handler http.Handler) http.Handler {
	headers := gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Internal", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With"})
	methods := gorillaHandlers.AllowedMethods([]string{"DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"})
	origins := gorillaHandlers.AllowedOrigins([]string{"*"})
	return gorillaHandlers.CORS(headers, methods, origins)(handler)
}

func logServerInfo(utilsContainer *utils.Container, servicesContainer *services.Container) {
	port := getServerPort(utilsContainer)
	serverURL := fmt.Sprintf("Server url: http://localhost%s", port)
	swaggerURL := fmt.Sprintf("%s/api-docs/index.html", serverURL)
	servicesContainer.Logger.Info("Initializing server")
	servicesContainer.Logger.Info(fmt.Sprintf("Environment: %s", utilsContainer.Environment.GetEnvVar("ENV")))
	servicesContainer.Logger.Info(serverURL)
	servicesContainer.Logger.Info(fmt.Sprintf("Swagger url: %s", swaggerURL))
}

func getServerPort(utilsContainer *utils.Container) string {
	port := fmt.Sprintf(":%v", utilsContainer.Environment.GetEnvVar("PORT"))
	if port == ":" {
		port = ":8080"
	}
	return port
}
