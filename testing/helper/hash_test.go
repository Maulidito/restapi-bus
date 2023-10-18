package helper

import (
	"fmt"
	"restapi-bus/helper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	pass := "CHECK"
	salt := fmt.Sprint(time.Now().UnixNano())

	dataEncrypt1 := helper.HashPasswordBcrypt(pass, salt)

	time.Sleep(time.Second * 2)

	dataEncrypt2 := helper.HashPasswordBcrypt(pass, salt)

	t.Log("pass =", pass)
	t.Log("Data Encrypt =", dataEncrypt1)
	t.Log("Data Encrypt =", dataEncrypt2)

	err1 := bcrypt.CompareHashAndPassword([]byte(dataEncrypt1), []byte(pass))
	err2 := bcrypt.CompareHashAndPassword([]byte(dataEncrypt2), []byte(pass))
	assert.Nil(t, err1)
	assert.Nil(t, err2)

}
