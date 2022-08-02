package integration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"restapi-bus/depedency"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var serverAddr = "http://localhost:8080/v1"
var server *gin.Engine

func TestMain(m *testing.M) {

	server = depedency.InitializedServer()

	extit := m.Run()

	os.Exit(extit)

}

func TestGetAllTicket(t *testing.T) {
	t.Log("GET ALL TICKET")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, serverAddr+"/ticket/", nil)

	server.ServeHTTP(recorder, request)

	// bodyByte, err := io.ReadAll(recorder.Result().Body)

	// helper.PanicIfError(err)

	assert.Equal(t, recorder.Code, http.StatusOK)

}
func TestGetOneTicket(t *testing.T) {
	id := "2"
	t.Log("GET TICKET WHERE TICKET ID", id)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodGet, serverAddr+"/ticket/"+id, nil)

	server.ServeHTTP(recorder, request)

	// bodyByte, err := io.ReadAll(recorder.Result().Body)

	// helper.PanicIfError(err)

	assert.Equal(t, recorder.Code, http.StatusOK)

}

func TestGetOneTicketFailed(t *testing.T) {
	id := "29"
	t.Log("GET TICKET WHERE TICKET ID", id)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodGet, serverAddr+"/ticket/"+id, nil)

	server.ServeHTTP(recorder, request)

	// bodyByte, err := io.ReadAll(recorder.Result().Body)

	// helper.PanicIfError(err)

	assert.Equal(t, recorder.Code, http.StatusNotFound)

}

func TestGetOneTicketFailedNotInt(t *testing.T) {
	id := "bukanInt"
	t.Log("GET TICKET WHERE TICKET ID", id)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodGet, serverAddr+"/ticket/"+id, nil)

	server.ServeHTTP(recorder, request)

	// bodyByte, err := io.ReadAll(recorder.Result().Body)

	// helper.PanicIfError(err)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)

}
