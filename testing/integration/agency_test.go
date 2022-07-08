package integration

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"restapi-bus/depedency"
	"restapi-bus/helper"
	"testing"

	"github.com/gin-gonic/gin"
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

	bodyByte, err := io.ReadAll(recorder.Result().Body)

	helper.PanicIfError(err)

	fmt.Println(string(bodyByte))

}
func TestGetOneTicket(t *testing.T) {
	id := "2"
	t.Log("GET TICKET WHERE TICKET ID", id)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodGet, serverAddr+"/ticket/"+id, nil)

	server.ServeHTTP(recorder, request)

	bodyByte, err := io.ReadAll(recorder.Result().Body)

	helper.PanicIfError(err)

	fmt.Println(string(bodyByte))

}
