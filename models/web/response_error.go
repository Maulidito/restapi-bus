package web

type ErrorMessage struct {
	Message string
}

type ResponseError struct {
	Code   int          `json:"code"`
	Status string       `json:"status"`
	Data   ErrorMessage `json:"data"`
}
