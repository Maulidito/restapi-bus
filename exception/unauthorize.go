package exception

type Unauthorized struct {
	Message string
}

func NewUnauthorizedError(errMsg string) Unauthorized {
	return Unauthorized{Message: errMsg}
}

func (unauthorized Unauthorized) Error() string {
	return unauthorized.Message
}
