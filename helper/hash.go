package helper

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashPassword(pass string, salt string) string {
	sha := sha256.New()
	_, err := sha.Write([]byte(pass + salt))
	PanicIfError(err)
	dataEncrypted := sha.Sum(nil)
	dataEncryptedEncode := base64.StdEncoding.EncodeToString(dataEncrypted)
	return dataEncryptedEncode
}
