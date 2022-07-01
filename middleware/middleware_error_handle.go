package middleware

import (
	"fmt"
	"net/http"
	"restapi-bus/models/web"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func MiddlewarePanicHandler(c *gin.Context) {
	defer func() {

		panicErr := recover()

		validatorErr, isValidator := panicErr.(validator.ValidationErrors)
		if isValidator {

			validatorErrorHandle(c, validatorErr)
		}
	}()

	c.Next()

}

func validatorErrorHandle(c *gin.Context, err validator.ValidationErrors) {

	FieldErrMessage := []web.ErrorMessage{}

	for _, v := range err {
		message := web.ErrorMessage{Message: fmt.Sprintf("ERROR BINDING DATA, WHAT = %s, WHERE = %s", v.ActualTag(), v.Field())}
		FieldErrMessage = append(FieldErrMessage, message)
	}

	response := web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: FieldErrMessage}
	c.JSON(http.StatusBadRequest, response)
}
