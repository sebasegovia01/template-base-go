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

	// init containers
	utilsContainer := newUtilsContainer()

	servicesContainer := newServicesContainer()

	handlersContainers := newHandlersContainer(servicesContainer)

	// Verifica si la aplicación se está ejecutando en AWS Lambda
	_, isLambda := utilsContainer.Environment.LookupEnv("AWS_LAMBDA_RUNTIME_API")
	fmt.Println("IsLambda?: ", isLambda)
	if isLambda {
		// Si es así, utiliza la función LambdaHandler como el manejador de la función Lambda
		lambda.Start(LambdaHandler)
	} else {
		// environment arg
		args := os.Args[1:]

		var localEnv string
		if len(args) > 0 {
			localEnv = os.Args[1:][0] // just 1 arg is received
		}

		// parse from .env file if exists
		if err := utils.ParseEnvironmentFile(localEnv); err != nil {
			log.Fatal(err)
		}

		// CORS
		headers := gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Internal", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With"})
		methods := gorillaHandlers.AllowedMethods([]string{"DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"})
		origins := gorillaHandlers.AllowedOrigins([]string{"*"})

		api := core.NewApi(handlersContainers, &utils.Logger{})

		servicesContainer.Logger.Info("Initializing server")

		port := fmt.Sprintf(":%v", utilsContainer.Environment.GetEnvVar("PORT"))

		if port == ":" {
			port = ":8080"
		}

		env := utilsContainer.Environment.GetEnvVar("ENV")

		if env == "" {
			env = "development"
		}
		servicesContainer.Logger.Info(fmt.Sprintf("Environment: %s", env))

		serverURL := fmt.Sprintf("Server url: http://localhost%s", port)
		servicesContainer.Logger.Info(serverURL)

		swaggerURL := fmt.Sprintf("%s/api-docs/index.html", serverURL)
		servicesContainer.Logger.Info(fmt.Sprintf("Swagger url: %s", swaggerURL))

		log.Fatal(http.ListenAndServe(port, gorillaHandlers.CORS(headers, methods, origins)(api.Router())))
	}
}

// function to handler LambdaFunctionURLRequest
// SE DEBE MODIFICAR EL EVENTO PARA HACER MATCH CON LA LLAMADA ESPECIFICA.
func LambdaHandler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// Inicializar contenedores como antes
	servicesContainer := newServicesContainer()
	handlersContainers := newHandlersContainer(servicesContainer)
	api := core.NewApi(handlersContainers, &utils.Logger{})

	// Convertir el evento de LambdaFunctionURL a una solicitud HTTP utilizando httptest.NewRecorder para capturar la respuesta
	req, err := utils.LambdaEventToHttpRequest(request)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}
	recorder := httptest.NewRecorder()

	// Usar el enrutador para manejar la solicitud convertida
	api.Router().ServeHTTP(recorder, req)

	// Extraer la respuesta del recorder y convertirla a una respuesta de LambdaFunctionURL
	response := events.LambdaFunctionURLResponse{
		StatusCode: recorder.Code,
		Body:       recorder.Body.String(),
		Headers:    map[string]string{"Content-Type": "application/json"}, // Ajustar según sea necesario
	}

	return response, nil
}
