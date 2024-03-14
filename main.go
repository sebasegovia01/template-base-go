package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"template-base-go/src/core"
	"template-base-go/src/handlers"
	"template-base-go/src/models"
	"template-base-go/src/services"
	"template-base-go/src/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	gorillaHandlers "github.com/gorilla/handlers"
)

var utilsContainer *utils.Container
var handlersContainer *handlers.Container
var servicesContainer *services.Container

func main() {
	loadEnvs()

	utilsContainer = utils.NewContainer(utils.NewEnvironment())

	isLambda := utilsContainer.Environment.IsLambda()

	fmt.Println("IsLambda?: ", isLambda)

	servicesContainer = services.NewContainer(&utils.Logger{}, &utils.Environment{}, &services.OtpService{})
	handlersContainer = handlers.NewContainer(handlers.NewOtpHandler(servicesContainer.Logger, servicesContainer.OtpService))

	if isLambda {
		lambda.Start(LambdaHandler)
	} else {
		httpServer()
	}

}

func loadEnvs() {
	if err := utils.ParseEnvironmentFile(""); err != nil {
		log.Fatal(err)
	}
}

func httpServer() {

	api := core.NewApi(handlersContainer, &utils.Logger{})

	port := getServerPort()

	logServerInfo()

	log.Fatal(http.ListenAndServe(port, setupCORS(api.Router())))
}

func LambdaHandler(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	// Imprime el evento recibido para depuración
	eventJSON, _ := json.Marshal(event)
	fmt.Printf("Received event: %s\n", string(eventJSON))

	if _, ok := event["requestContext"]; ok {
		fmt.Println("Identified as Lambda Function URL Request")

		lambdaRequest, err := utils.ConvertToLambdaFunctionURLRequest(event)
		if err != nil {
			return nil, fmt.Errorf("failed to convert event to LambdaFunctionURLRequest: %v", err)
		}
		return handleLambdaFunctionURLRequest(lambdaRequest)
	}

	// Step Function Event
	fmt.Println("Handling as Step Function Event or other JSON event")

	rawMessage, err := utils.ConvertToJSONRawMessage(event)
	if err != nil {
		return nil, fmt.Errorf("failed to convert event to JSONRawMessage: %v", err)
	}

	return handleJSONRawMessageEvent(rawMessage)
}

func handleJSONRawMessageEvent(rawMessage json.RawMessage) (interface{}, error) {
	// Lógica para manejar el evento
	var payload models.SFPayload
	if err := json.Unmarshal(rawMessage, &payload); err != nil {
		fmt.Printf("Error unmarshaling JSON Raw Message Event: %v\n", err)
		return nil, err
	}

	fmt.Printf("Handling JSON Raw Message Payload: %+v\n", payload)

	// Crear una solicitud HTTP falsa para simular la llamada a la API
	req, err := http.NewRequest(http.MethodGet, payload.Path, nil)
	if err != nil {
		fmt.Printf("Error creating fake HTTP request: %v\n", err)
		return nil, err
	}

	// Inicializar la API y manejar la solicitud HTTP falsa
	api := core.NewApi(handlersContainer, &utils.Logger{})
	recorder := httptest.NewRecorder()
	api.Router().ServeHTTP(recorder, req)

	// Devolver la respuesta HTTP como una respuesta de Step Function
	return map[string]interface{}{
		"statusCode": recorder.Code,
		"body":       recorder.Body.String(),
		"headers":    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func handleLambdaFunctionURLRequest(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	api := core.NewApi(handlersContainer, &utils.Logger{})
	req, err := utils.LambdaEventToHttpRequest(request)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}, nil
	}

	recorder := httptest.NewRecorder()
	api.Router().ServeHTTP(recorder, req)

	return events.LambdaFunctionURLResponse{
		StatusCode: recorder.Code,
		Body:       recorder.Body.String(),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func setupCORS(handler http.Handler) http.Handler {
	headers := gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With"})
	methods := gorillaHandlers.AllowedMethods([]string{"DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"})
	origins := gorillaHandlers.AllowedOrigins([]string{"*"})
	return gorillaHandlers.CORS(headers, methods, origins)(handler)
}

func logServerInfo() {
	port := getServerPort()
	serverURL := fmt.Sprintf("Server url: http://localhost%s", port)
	swaggerURL := fmt.Sprintf("%s/api-docs/index.html", serverURL)
	servicesContainer.Logger.Info("Initializing server")
	servicesContainer.Logger.Info(fmt.Sprintf("Environment: %s", utilsContainer.Environment.GetEnvVar("ENV")))
	servicesContainer.Logger.Info(serverURL)
	servicesContainer.Logger.Info(fmt.Sprintf("Swagger url: %s", swaggerURL))
}

func getServerPort() string {
	port := fmt.Sprintf(":%v", utilsContainer.Environment.GetEnvVar("PORT"))
	if port == ":" {
		port = ":8080"
	}
	return port
}
