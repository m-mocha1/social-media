package main

import "golang.org/x/crypto/bcrypt"

func passwordHash(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}

func checkPass(p, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
	return err == nil // return a true of false
}
