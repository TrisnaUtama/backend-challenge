package health

type Service interface {
	Ping() PingResponse
	Echo(req EchoRequest) EchoResponse
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Ping() PingResponse {
	return PingResponse{
		Message: "pong",
	}
}

func (s *service) Echo(req EchoRequest) EchoResponse {
	return EchoResponse{
		Message: req.Message,
	}
}
