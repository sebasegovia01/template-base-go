package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"template-base-go/src/core"
	"template-base-go/src/handlers"
	"template-base-go/src/repositories"
	"template-base-go/src/services"
	"template-base-go/src/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	gorillaHandlers "github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

var utilsContainer *utils.Container
var handlersContainer *handlers.Container
var servicesContainer *services.Container

func main() {
	loadEnvs()

	utilsContainer = utils.NewContainer(utils.NewEnvironment())

	isLambda := utilsContainer.Environment.IsLambda()

	fmt.Println("IsLambda?: ", isLambda)

	servicesContainer = services.NewContainer(&utils.Logger{}, &utils.Environment{}, &services.ExampleService{})
	dbConnection := getDbConnection()
	repositoriesContainer := repositories.NewContainer(repositories.NewExampleRepository(dbConnection))
	handlersContainer = handlers.NewContainer(handlers.NewExampleHandler(servicesContainer.Logger, servicesContainer.ExampleServices, repositoriesContainer.ExampleRepository))

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

func LambdaHandler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	api := core.NewApi(handlersContainer, &utils.Logger{})

	logServerInfo()

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

func getDbConnection() *mongo.Database {

	var (
		uri     = utilsContainer.Environment.GetEnvVar("MONGO_URI")
		db_name = utilsContainer.Environment.GetEnvVar("DB_NAME")
	)

	db, err := repositories.Connect(uri, db_name)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
