package repositories

type Container struct {
	OTPRepository IOTPRepository
}

func NewContainer(
	OTPRepository IOTPRepository,
) *Container {
	return &Container{
		OTPRepository: OTPRepository,
	}
}
