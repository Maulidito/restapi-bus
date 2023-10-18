package helper

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string, salt string) string {
	sha := sha256.New()
	_, err := sha.Write([]byte(pass + salt))
	PanicIfError(err)
	dataEncrypted := sha.Sum(nil)
	dataEncryptedEncode := base64.StdEncoding.EncodeToString(dataEncrypted)
	return dataEncryptedEncode
}

func HashPasswordBcrypt(pass string, salt string) string {
	dataGenerate, err := bcrypt.GenerateFromPassword([]byte(pass+salt), bcrypt.DefaultCost)
	PanicIfError(err)

	return string(dataGenerate)
}

func HashPasswordCompare(hashedPassword string, pass string, salt string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass+salt))
	return err
}
