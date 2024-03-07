package handlers

import (
	"encoding/json"
	"net/http"
	"template-base-go/src/services"
	"template-base-go/src/utils"
)

type IExampleHandler interface {
	Do(w http.ResponseWriter, r *http.Request)
}

type ExampleHandler struct {
	logger         utils.ILogger
	exampleService services.IExampleService
}

// Container
func NewExampleHandler(
	logger utils.ILogger,
	exampleService services.IExampleService,
) *ExampleHandler {
	return &ExampleHandler{
		logger,
		exampleService,
	}
}

func (h *ExampleHandler) Do(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Calling Do example handler")

	exampleResponse := struct {
		Message string `json:"message"`
	}{
		Message: h.exampleService.SayHelloTo("John Doe"),
	}

	w.WriteHeader(http.StatusOK)

	// Codifica el objeto de respuesta y lo escribe (y envia) en el ResponseWriter
	if err := json.NewEncoder(w).Encode(exampleResponse); err != nil {
		// En caso de error al codificar la respuesta, se registra el error.
		// La respuesta al cliente ya fue enviada, así que este error sería principalmente
		// para propósitos de logging o seguimiento de fallos.
		h.logger.Error("Failed to encode response", map[string]interface{}{"error": err})
	}
}
