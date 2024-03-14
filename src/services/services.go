package services

import (
	"template-base-go/src/utils"
)

type Container struct {
	Logger      utils.ILogger
	Environment utils.IEnvironment
	OtpService  IOtpService
}

func NewContainer(
	logger utils.ILogger,
	env utils.IEnvironment,
	otpService IOtpService,
) *Container {
	return &Container{
		logger,
		env,
		otpService,
	}
}
