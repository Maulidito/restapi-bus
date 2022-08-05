package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"restapi-bus/constant"
	"restapi-bus/models/web"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JwtBearer struct {
	BearerToken string `header:"Authorization" binding:"required,startswith=Bearer"`
}

func MiddlewareAuth(ctx *gin.Context) {

	JwtBearer := JwtBearer{}
	err := ctx.BindHeader(&JwtBearer)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.ResponseError{Code: http.StatusUnauthorized, Status: "UNAUTHORIZED", Data: web.ErrorMessage{ErrorMessage: err.Error()}})
		return

	}

	Token := JwtBearer.BearerToken

	if Token == "" {
		Token, err = ctx.Cookie(constant.X_API_KEY)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.ResponseError{Code: http.StatusUnauthorized, Status: "UNAUTHORIZED", Data: web.ErrorMessage{ErrorMessage: err.Error()}})
			return
		}
	}
	onlyToken := strings.Replace(Token, "Bearer ", "", 1)
	claim := &jwt.RegisteredClaims{}

	tokenAfterParse, err := jwt.ParseWithClaims(onlyToken, claim, func(t *jwt.Token) (interface{}, error) {

		err := t.Claims.Valid()
		if err != nil {
			return nil, err
		}
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("error signing method")
		}
		claim.ExpiresAt = jwt.NewNumericDate(claim.ExpiresAt.Add(time.Minute * 1))
		secret := os.Getenv("SECRET_KEY_AUTH")
		return []byte(secret), nil
	})

	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); !ok {
			ctx.AbortWithStatusJSON(int(verr.Errors), web.ResponseError{Code: int(verr.Errors), Status: "UNAUTHORIZED", Data: web.ErrorMessage{ErrorMessage: verr.Inner.Error()}})
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.ResponseError{Code: http.StatusUnauthorized, Status: "UNAUTHORIZED", Data: web.ErrorMessage{ErrorMessage: err.Error()}})
		return
	}

	if !tokenAfterParse.Valid {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.ResponseError{Code: http.StatusUnauthorized, Status: "UNAUTHORIZED", Data: web.ErrorMessage{ErrorMessage: err.Error()}})
		return
	}
	ctx.Next()

}
