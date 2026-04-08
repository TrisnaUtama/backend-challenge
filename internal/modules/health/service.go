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
		Success: true,
	}
}

func (s *service) Echo(req EchoRequest) EchoResponse {
	return EchoResponse{
		Success: true,
	}
}
