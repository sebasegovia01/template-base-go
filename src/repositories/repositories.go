package repositories

type Container struct {
	ExampleRepository IExampleRepository
}

func NewContainer(
	exampleRepository IExampleRepository,
) *Container {
	return &Container{
		ExampleRepository: exampleRepository,
	}
}
