package handlers

import (
	"encoding/json"
	"net/http"
	"put-otp-go/src/models"
	"put-otp-go/src/repositories"
	"put-otp-go/src/utils"

	"github.com/gorilla/mux"
)

type IOtpHandler interface {
	Do(w http.ResponseWriter, r *http.Request)
}

type OtpHandler struct {
	logger        utils.ILogger
	oTPRepository repositories.IOTPRepository
}

// Container
func NewOtpHandler(
	logger utils.ILogger,
	oTPRepository repositories.IOTPRepository,
) *OtpHandler {
	return &OtpHandler{
		logger,
		oTPRepository,
	}
}

// UpdateOtp godoc
// @Summary Update an OTP
// @Description update otp by id and body data
// @Tags Update OTP
// @Accept json
// @Produce json
// @Param id path string true "OTP ID"
// @Param otp body models.Otp true "OTP object"
// @Success 200 {object} models.UpdateResponse
// @Router /service-otp/v1/put/{id} [put]
func (h *OtpHandler) Update(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Calling Update otp handler")

	// Get id from params
	vars := mux.Vars(r)
	id := vars["id"]

	var otp models.Otp
	if err := json.NewDecoder(r.Body).Decode(&otp); err != nil {
		h.logger.Error("Failed to decode request body", map[string]interface{}{"error": err})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.oTPRepository.Update(id, otp)
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
