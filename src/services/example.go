package services

import "template-base-go/src/utils"

type IExampleService interface {
	SayHelloTo(name string) string
}

type ExampleService struct {
	logger utils.ILogger
}

func NewService(logger utils.ILogger) ExampleService {
	return ExampleService{
		logger: logger,
	}
}

func (s *ExampleService) SayHelloTo(name string) string {

	return "Hello " + name

}
