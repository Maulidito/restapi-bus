package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"restapi-bus/constant"
	"restapi-bus/exception"
	"restapi-bus/helper"
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
	helper.PanicIfError(err)

	Token := JwtBearer.BearerToken

	if Token == "" {
		Token, err = ctx.Cookie(constant.X_API_KEY)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.ResponseError{Code: http.StatusUnauthorized, Status: "UNAUTHORIZED", Data: web.ErrorMessage{ErrorMessage: err.Error()}})
			return
		}
	}
	onlyToken := strings.Replace(Token, "Bearer ", "", 1)
	claim := &web.Claim{}

	tokenAfterParse, err := jwt.ParseWithClaims(onlyToken, claim, func(t *jwt.Token) (interface{}, error) {

		err := t.Claims.Valid()
		if err != nil {
			return nil, err
		}
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("error signing method")
		}

		secret := os.Getenv("SECRET_KEY_AUTH")
		return []byte(secret), nil
	})

	if err != nil {
		ctx.Abort()
		panic(exception.NewUnauthorizedError(err.Error()))

	}

	if !tokenAfterParse.Valid {
		ctx.Abort()
		panic(exception.NewUnauthorizedError("token invalid"))

	}
	fmt.Println(claim)
	if newToken := refreshToken(claim, tokenAfterParse); newToken != "" {
		ctx.SetCookie(constant.X_API_KEY, newToken, claim.ExpiresAt.Second(), "/", "localhost", true, true)
	}
	ctx.Next()

}

//check if token expired token less than 1 hour than refresh
func refreshToken(claim *web.Claim, token *jwt.Token) string {
	timeToRefresh := claim.ExpiresAt.Time

	if time.Now().After(timeToRefresh.Add(-time.Hour * 1)) {
		claim.ExpiresAt = jwt.NewNumericDate(claim.ExpiresAt.Add(time.Hour * 1))
	}
	tokenstring, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY_AUTH")))
	return tokenstring

}
