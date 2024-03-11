package handlers

import (
	"encoding/json"
	"net/http"
	"template-base-go/src/models"
	"template-base-go/src/repositories"
	"template-base-go/src/services"
	"template-base-go/src/utils"

	"github.com/gorilla/mux"
)

type IExampleHandler interface {
	Do(w http.ResponseWriter, r *http.Request)
}

type ExampleHandler struct {
	logger            utils.ILogger
	exampleService    services.IExampleService
	exampleRepository repositories.IExampleRepository
}

// Container
func NewExampleHandler(
	logger utils.ILogger,
	exampleService services.IExampleService,
	exampleRepository repositories.IExampleRepository,
) *ExampleHandler {
	return &ExampleHandler{
		logger,
		exampleService,
		exampleRepository,
	}
}

func (h *ExampleHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Calling Get example handler")

	vars := mux.Vars(r)
	id := vars["id"]

	response, err := h.exampleRepository.Get(id)

	if err != nil {
		h.logger.Error(err.Error())
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response", map[string]interface{}{"error": err})
	}
}

func (h *ExampleHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Calling Create example handler")

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.logger.Error("Failed to decode request body", map[string]interface{}{"error": err})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.exampleRepository.Create(book)
	if err != nil {
		h.logger.Error("Failed to create book", map[string]interface{}{"error": err})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error("Failed to encode response", map[string]interface{}{"error": err})
	}
}

func (h *ExampleHandler) Update(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Calling Update example handler")

	// Get id from params
	vars := mux.Vars(r)
	id := vars["id"]

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.logger.Error("Failed to decode request body", map[string]interface{}{"error": err})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.exampleRepository.Update(id, book)
	if err != nil {
		h.logger.Error("Failed to update book", map[string]interface{}{"error": err})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error("Failed to encode response", map[string]interface{}{"error": err})
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

	if err := json.NewEncoder(w).Encode(exampleResponse); err != nil {
		h.logger.Error("Failed to encode response", map[string]interface{}{"error": err})
	}
}
