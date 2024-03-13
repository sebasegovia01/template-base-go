package utils

import (
	"bufio"
	"bytes"
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
	Logger      ILogger
}

func NewContainer(logger ILogger, env IEnvironment) *Container {
	return &Container{
		Environment: env,
		Logger:      logger,
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

// LambdaEventToHttpRequest convierte un evento en una solicitud HTTP
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
