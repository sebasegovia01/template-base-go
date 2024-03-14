package handlers

type Container struct {
	//handlers *definitions here
	*OtpHandler
}

func NewContainer(otpHandler *OtpHandler) *Container {
	return &Container{
		//handlers instances here
		OtpHandler: otpHandler,
	}
}
