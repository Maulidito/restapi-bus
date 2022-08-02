package integration

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"restapi-bus/helper"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAgency(t *testing.T) {

	t.Log("Add Ticket")
	agencyName := "AgencyName"
	agencyPlace := "AgencyPlace"
	agencyUsername := "AgencyUsername@gmail.com"
	agencyPassword := "agendfsfs"
	bodyRequest := strings.NewReader(fmt.Sprintf(`name=%s&place=%s&username=%s&password=%s`, agencyName, agencyPlace, agencyUsername, agencyPassword))

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodPost, serverAddr+"/agency/", bodyRequest)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	server.ServeHTTP(recorder, request)

	bodyByte, err := io.ReadAll(recorder.Result().Body)

	helper.PanicIfError(err)

	t.Log(request.Form.Get("password"))

	assert.Equal(t, recorder.Code, http.StatusOK, string(bodyByte))

}
