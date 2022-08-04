package web

import "github.com/golang-jwt/jwt/v4"

type Claim struct {
	*jwt.RegisteredClaims
	Username string
}

type Token struct {
	Token string `json:"token"`
}

type WebResponseToken struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   *Token `json:"data"`
}
