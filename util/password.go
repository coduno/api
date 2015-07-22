package util

import (
  "crypto/rand"
  "golang.org/x/crypto/bcrypt"
)


const passwordLength = 16
const hashCost = 12

func GenerateRandomPassword() string{
  pw := make([]byte, passwordLength)
	rand.Read(pw)
	return string(pw)
}

func EncryptPassword(password string) string{
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
  if err != nil {
    panic(err)
  }
  return string(hashedPassword)
}

func CheckPassword(hashedPassword, password string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
  if err == nil{
    return true
  }
  return false
}
