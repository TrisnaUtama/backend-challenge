package health

type PingResponse struct {
	Success bool `json:"message"`
}

type EchoRequest struct {
	Message string `json:"message"`
}

type EchoResponse struct {
	Success bool `json:"message"`
}
