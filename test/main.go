package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	fmt.Println(bytes)
}
