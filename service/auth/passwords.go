package auth

import (
	"fmt"
	"github.com/Code-xCartel/noxus-api-svc/config"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
}

func GenerateUniqueId(length int) string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	chars := config.Envs.CharStrings
	rand.Seed(time.Now().UnixNano())
	var randomString strings.Builder
	for i := 0; i < length; i++ {
		randomString.WriteByte(chars[rand.Intn(len(chars))])
	}
	id := fmt.Sprintf("%d%s", timestamp, randomString.String())
	return fmt.Sprintf("NOX-%s", id[:length])
}
