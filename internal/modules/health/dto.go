package health

type PingResponse struct {
	Success bool `json:"success"`
}

type EchoRequest struct {
	Message string `json:"message"`
}

type EchoResponse struct {
	Success bool `json:"success"`
}
