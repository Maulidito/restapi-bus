package exception

type NotFound struct {
	Message string
}

func NewNotFoundError(errMsg string) NotFound {
	return NotFound{Message: errMsg}
}
