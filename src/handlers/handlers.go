package handlers

type Container struct {
	//handlers *definitions here
	*ExampleHandler
}

func NewContainer(exampleHandler *ExampleHandler) *Container {
	return &Container{
		//handlers instances here
		ExampleHandler: exampleHandler,
	}
}
