package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Container struct {
	Environment IEnvironment
}

func NewContainer(env IEnvironment) *Container {
	return &Container{
		Environment: env,
	}
}

// ParseEnvironmentFile checks if running in local environment and reads .env file if present
func ParseEnvironmentFile(localEnv string) error {
	// Check if running in a local environment (e.g., not in AWS Lambda)
	if _, isLambda := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); !isLambda {
		file := ".env"
		if localEnv != "" {
			file = fmt.Sprintf(".env%s", localEnv)
		}

		handler, err := os.Open(file)
		if err != nil {
			log.Printf("Failed to open file path, %v", err)
			return nil // Don't fail in Lambda if .env file is not found
		}
		defer handler.Close()

		return ReadFileAndSetEnv(handler)
	}
	return nil
}

// ReadFileAndSetEnv takes a reader and sets its keys as environment variables
func ReadFileAndSetEnv(handle io.Reader) error {
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Ignoring malformed line in env file: %s", line)
			continue
		}
		key := parts[0]
		value := parts[1]
		os.Setenv(key, value)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading env file: %v", err)
		return err
	}
	return nil
}

// LambdaResponseWriter implementa http.ResponseWriter para capturar la respuesta HTTP
type LambdaResponseWriter struct {
	StatusCode int
	Body       bytes.Buffer
	Headers    http.Header
}

// LambdaEventToHttpRequest convierte un evento APIGatewayProxyRequest en una solicitud HTTP
func LambdaEventToHttpRequest(req events.LambdaFunctionURLRequest) (*http.Request, error) {
	// Convertir el cuerpo del evento Lambda en un io.Reader
	reader := bytes.NewBufferString(req.Body)

	// Crear la solicitud HTTP
	httpReq, err := http.NewRequest(strings.ToUpper(req.RequestContext.HTTP.Method), req.RequestContext.HTTP.Path, reader)
	if err != nil {
		return nil, err
	}

	// Copiar los encabezados al HTTP Request
	for key, value := range req.Headers {
		httpReq.Header.Add(key, value)
	}

	return httpReq, nil
}

// NewLambdaResponseWriter crea una nueva instancia de LambdaResponseWriter
func NewLambdaResponseWriter() *LambdaResponseWriter {
	return &LambdaResponseWriter{
		Headers: http.Header{},
	}
}

func (lrw *LambdaResponseWriter) Header() http.Header {
	return lrw.Headers
}

func (lrw *LambdaResponseWriter) Write(data []byte) (int, error) {
	return lrw.Body.Write(data)
}

func (lrw *LambdaResponseWriter) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
}

// ToLambdaResponse convierte la respuesta capturada en una respuesta de API Gateway
func (lrw *LambdaResponseWriter) ToLambdaResponse() (events.APIGatewayProxyResponse, error) {
	var body string
	if lrw.Body.Len() > 0 {
		body = lrw.Body.String()
	} else {
		// Si el cuerpo está vacío, intenta devolver un JSON vacío para mantener la consistencia
		emptyJSON, _ := json.Marshal(map[string]interface{}{})
		body = string(emptyJSON)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: lrw.StatusCode,
		Body:       body,
		Headers:    convertHeaders(lrw.Headers),
	}, nil
}

// convertHeaders convierte http.Header a map[string]string requerido por APIGatewayProxyResponse
func convertHeaders(h http.Header) map[string]string {
	headers := map[string]string{}
	for k, v := range h {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	return headers
}
