package middleware

import (
	"fmt"
	"net/http"
	"restapi-bus/models/web"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func MiddlewarePanic(ctx *gin.Context) {
	defer func() {

		panicErr := recover()

		validatorErr, isValidator := panicErr.(validator.ValidationErrors)
		if isValidator {

			validatorErrorHandle(ctx, validatorErr)
			return
		}
		NotFoundErrorHandle(ctx, panicErr)
	}()

	ctx.Next()

}

func validatorErrorHandle(c *gin.Context, err validator.ValidationErrors) {

	FieldErrMessage := []web.ErrorMessage{}

	for _, v := range err {
		message := web.ErrorMessage{ErrorMessage: fmt.Sprintf("ERROR BINDING DATA, WHAT = %s, WHERE = %s", v.ActualTag(), v.Field())}
		FieldErrMessage = append(FieldErrMessage, message)
	}

	response := web.ResponseBindingError{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: FieldErrMessage}
	c.JSON(http.StatusBadRequest, response)
}

func NotFoundErrorHandle(c *gin.Context, err interface{}) {

	dataErr, isString := err.(string)
	if !isString {
		return
	}
	errMsg := web.ErrorMessage{ErrorMessage: dataErr}
	response := web.ResponseError{Code: http.StatusNotFound, Status: "NOT FOUND", Data: errMsg}
	c.JSON(http.StatusNotFound, response)
}
