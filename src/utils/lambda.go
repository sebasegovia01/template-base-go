package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

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

// StepFunctionEventToHttpRequest convierte un evento de Step Function en una solicitud HTTP.
func StepFunctionEventToHttpRequest(event json.RawMessage) (*http.Request, error) {
	// Deserializar el evento de Step Function en un objeto Go.
	var eventData map[string]interface{}
	if err := json.Unmarshal(event, &eventData); err != nil {
		return nil, fmt.Errorf("error unmarshaling Step Function event: %v", err)
	}

	// Aquí puedes extraer y utilizar la información del evento para construir la solicitud HTTP.
	// En este ejemplo, asumimos que el evento contiene una ruta y un cuerpo para la solicitud HTTP.
	path, ok := eventData["path"].(string)
	if !ok {
		return nil, fmt.Errorf("event does not contain a valid 'path' field")
	}

	bodyData, ok := eventData["body"]
	if !ok {
		return nil, fmt.Errorf("event does not contain a 'body' field")
	}
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling body data: %v", err)
	}

	// Crear la solicitud HTTP.
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Establecer encabezados adicionales según sea necesario.
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func ConvertToLambdaFunctionURLRequest(event map[string]interface{}) (events.LambdaFunctionURLRequest, error) {
	// Convierte el mapa a json.RawMessage
	rawData, err := json.Marshal(event)
	if err != nil {
		return events.LambdaFunctionURLRequest{}, fmt.Errorf("error marshaling event to JSON: %v", err)
	}

	// Deserializa json.RawMessage en events.LambdaFunctionURLRequest
	var lambdaRequest events.LambdaFunctionURLRequest
	err = json.Unmarshal(rawData, &lambdaRequest)
	if err != nil {
		return events.LambdaFunctionURLRequest{}, fmt.Errorf("error unmarshaling JSON to LambdaFunctionURLRequest: %v", err)
	}

	return lambdaRequest, nil
}

func ConvertToJSONRawMessage(event map[string]interface{}) (json.RawMessage, error) {
	// Convierte el mapa a json.RawMessage
	rawData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("error marshaling event to JSON: %v", err)
	}

	return rawData, nil
}
