package middleware

import (
	"fmt"
	"net/http"
	"restapi-bus/models/web"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func MiddlewarePanic(ctx *gin.Context, panicErr any) {
	defer func() {

		if panicErr == nil {
			return
		}
		if validatorErrorHandle(ctx, &panicErr) {
			return
		}

		if NotFoundErrorHandle(ctx, panicErr) {
			return
		}

		if BadParamRouter(ctx, panicErr) {
			return
		}
	}()

	ctx.Next()

}

func validatorErrorHandle(c *gin.Context, err interface{}) bool {

	validatorErr, isValidator := err.(validator.ValidationErrors)
	if !isValidator {
		fmt.Println("OUT FROM VALIDATOR ERROR")
		return false
	}

	FieldErrMessage := []web.ErrorMessage{}

	for _, v := range validatorErr {
		message := web.ErrorMessage{ErrorMessage: fmt.Sprintf("ERROR BINDING DATA, WHAT = %s, WHERE = %s", v.ActualTag(), v.Field())}
		FieldErrMessage = append(FieldErrMessage, message)
	}

	response := web.ResponseBindingError{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: FieldErrMessage}
	c.JSON(http.StatusBadRequest, response)
	return true
}

func NotFoundErrorHandle(c *gin.Context, err interface{}) bool {

	dataErr, isString := err.(string)
	if !isString {
		fmt.Println(err)
		return false
	}
	errMsg := web.ErrorMessage{ErrorMessage: dataErr}
	response := web.ResponseError{Code: http.StatusNotFound, Status: "NOT FOUND", Data: errMsg}
	c.JSON(http.StatusNotFound, response)
	return true
}

func BadParamRouter(c *gin.Context, err interface{}) bool {
	dataErr, isErr := err.(error)

	if !isErr {
		return false
	}
	errMsg := web.ErrorMessage{ErrorMessage: dataErr.Error()}
	response := web.ResponseError{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: errMsg}
	c.JSON(http.StatusBadRequest, response)
	return true
}
