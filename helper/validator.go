package helper

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var ValidateFromToDate validator.Func = func(fl validator.FieldLevel) bool {

	date := fl.Field().String()

	fromDate, err := time.Parse("2006-01-02", date)
	PanicIfError(err)

	param := fl.Param()
	if param == "" {
		panic(errors.New("no Parameter on tag"))
	}

	toDate, err := time.Parse("2006-01-02", fl.Parent().FieldByName(param).String())

	PanicIfError(err)

	fmt.Println("CHECk CUSTOM VALIDATION", fromDate, "CHECK PARAM", toDate)

	return fromDate.Before(toDate)

}
