package services

import (
	"template-base-go/src/utils"
)

type Container struct {
	Logger          utils.ILogger
	Environment     utils.IEnvironment
	ExampleServices IExampleService
}

func NewContainer(
	logger utils.ILogger,
	env utils.IEnvironment,
	exampleService IExampleService,
) *Container {
	return &Container{
		logger,
		env,
		exampleService,
	}
}
