package helper

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

var ValidateFromToDate validator.Func = func(fl validator.FieldLevel) bool {
	param := fl.Param()
	if toDateString := fl.Parent().FieldByName(param).String(); toDateString != "" {

		date := fl.Field().String()

		fromDate, err := time.Parse("2006-01-02", date)
		PanicIfError(err)

		if param == "" {
			panic(errors.New("no Parameter on tag"))
		}

		toDate, err := time.Parse("2006-01-02", fl.Parent().FieldByName(param).String())

		PanicIfError(err)

		return fromDate.Before(toDate)
	}
	return true
}

var ValidateDateAfterNow validator.Func = func(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	if date == "" {
		return false
	}

	dateData, err := time.Parse("2006-01-02", date)
	PanicIfError(err)

	return dateData.After(time.Now())

}
