package web

type ErrorMessage struct {
	ErrorMessage string `json:"errorMessage"`
}

type ResponseError struct {
	Code   int          `json:"code"`
	Status string       `json:"status"`
	Data   ErrorMessage `json:"data"`
}

type ResponseBindingError struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Data   []ErrorMessage `json:"data"`
}
