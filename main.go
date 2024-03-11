package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"template-base-go/src/core"
	"template-base-go/src/handlers"
	"template-base-go/src/services"
	"template-base-go/src/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	_, isLambda := utilsContainer.Environment.LookupEnv("AWS_LAMBDA_RUNTIME_API")
	fmt.Println("IsLambda?: ", isLambda)

	if isLambda {
		lambda.Start(LambdaHandler)
	} else {
		startHTTPServer(utilsContainer, servicesContainer)
	}
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

func LambdaHandler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	api := setupAPI()

	req, err := utils.LambdaEventToHttpRequest(request)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	recorder := httptest.NewRecorder()
	api.Router().ServeHTTP(recorder, req)

	response := events.LambdaFunctionURLResponse{
		StatusCode: recorder.Code,
		Body:       recorder.Body.String(),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}

	return response, nil
}
