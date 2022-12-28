package utilities

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"math/rand"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Hash a plain password with the bcrypt library
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

// Check a hashed password with a plain password
func CheckPassword(hashedPassword string, plainPassword string) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPassword))

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Hash with SHA512
func HashSHA512(input string) string {
	// Get the bytes of the input
	inputBytes := []byte(input)
	// Create the SHA512 Hasher
	sha512obj := sha512.New()

	// Write input bytes to the hasher
	sha512obj.Write(inputBytes)

	// Get the hashed bytes
	hashedInputBytes := sha512obj.Sum(nil)

	// Convert hashed bytes to hex string
	hashedInputHex := hex.EncodeToString(hashedInputBytes)

	return hashedInputHex
}

// Check if SHA512 hashes match
func CheckSHA512(hashedInput string, input string) bool {
	hashedNewInput := HashPassword(input)
	return hashedNewInput == hashedInput
}

// Create a random string
func RandomStringWithUUID(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b) + (uuid.New().String())
}
