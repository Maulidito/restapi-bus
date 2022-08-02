package helper

import (
	"fmt"
	"restapi-bus/helper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	pass := "CHECK"
	salt := fmt.Sprint(time.Now().UnixNano())

	dataEncrypt1 := helper.HashPassword(pass, salt)

	time.Sleep(time.Second * 2)

	dataEncrypt2 := helper.HashPassword(pass, salt)

	t.Log("pass =", pass)
	t.Log("Data Encrypt =", dataEncrypt1)
	t.Log("Data Encrypt =", dataEncrypt2)

	assert.Equal(t, dataEncrypt1, dataEncrypt2)

}
