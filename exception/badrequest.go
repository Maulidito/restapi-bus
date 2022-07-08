package exception

type BadRequest struct {
	Message string
}

func NewBadRequestError(errMsg string) BadRequest {
	return BadRequest{Message: errMsg}
}
