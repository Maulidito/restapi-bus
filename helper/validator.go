package helper

import (
	"errors"
	"fmt"
	"strconv"
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

	return fromDate.Before(toDate)

}

var IsBool validator.Func = func(fl validator.FieldLevel) bool {
	fmt.Println("CHECK IN VALIDATOR IS BOOL")
	dataBool := fl.Field().String()

	_, err := strconv.ParseBool(dataBool)
	fmt.Println("CHECK DATA BOOL", fl.Field().String())
	if err != nil {
		return false
	}
	return true
}
