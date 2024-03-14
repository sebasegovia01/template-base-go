package services

import (
	"math/rand"
	"template-base-go/src/utils"
)

type IOtpService interface {
	GenerateRandomOtp(length int) string
}

type OtpService struct {
	logger utils.ILogger
}

func NewOtpService(logger utils.ILogger) OtpService {
	return OtpService{
		logger: logger,
	}
}

func (s *OtpService) GenerateRandomOtp(length int) string {

	if length <= 0 {
		length = 6 // pred value
	}

	digits := "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}

	return string(otp)

}
