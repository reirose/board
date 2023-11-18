package src

import (
	"crypto/rand"
	"errors"
	"fmt"
)

func generateRandomToken(len int) (string, error) {
	if len <= 0 {
		return "", errors.New("token generation error: insufficient length")
	}
	if len > 32 {
		len = 32
	}

	b := make([]byte, len*8)

	_, err := rand.Read(b)
	Catch(err)

	return fmt.Sprintf("%x", b), nil
}

func GenerateToken(userId string) string {
	token, err := generateRandomToken(len(userId))
	Catch(err)

	return userId + ":" + token
}
