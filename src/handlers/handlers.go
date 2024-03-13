package handlers

type Container struct {
	*OtpHandler
}

func NewContainer(otpHandler *OtpHandler) *Container {
	return &Container{
		OtpHandler: otpHandler,
	}
}
