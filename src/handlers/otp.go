package handlers

import (
	"encoding/json"
	"net/http"
	"template-base-go/src/models"
	"template-base-go/src/services"
	"template-base-go/src/utils"
)

type IOtpHandler interface {
	Generate(w http.ResponseWriter, r *http.Request)
}

type OtpHandler struct {
	logger     utils.ILogger
	otpService services.IOtpService
}

// Container
func NewOtpHandler(
	logger utils.ILogger,
	otpService services.IOtpService,
) *OtpHandler {
	return &OtpHandler{
		logger,
		otpService,
	}
}

// Generate godoc
// @Summary Generate OTP
// @Description Generates a random OTP of length 6
// @Tags OTP
// @Accept json
// @Produce json
// @Success 200 {object} models.Otp
// @Router /otp/generate [get]
func (h *OtpHandler) Generate(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Calling Generate otp handler")

	otp := models.Otp{
		Otp: h.otpService.GenerateRandomOtp(6),
	}

	h.logger.Info("OTP generated succesfully: " + otp.Otp)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(otp); err != nil {
		h.logger.Error("Failed to encode response", map[string]interface{}{"error": err})
	}
}
