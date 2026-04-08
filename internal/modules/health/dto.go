package health

type PingResponse struct {
	Message string `json:"message"`
}

type EchoRequest struct {
	Message string `json:"message"`
}

type EchoResponse struct {
	Message string `json:"message"`
}
