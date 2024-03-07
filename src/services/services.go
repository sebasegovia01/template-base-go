package services

import "template-base-go/src/utils"

type Container struct {
	Logger      utils.ILogger
	Environment utils.IEnvironment
}

// service container, add your services for reference inyection
func NewContainer(
	logger utils.ILogger,
	env utils.IEnvironment,
) *Container {
	return &Container{
		logger,
		env,
	}
}
