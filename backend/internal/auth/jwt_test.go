package auth

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestValidateToken(t *testing.T) {
	godotenv.Load("../../.env.test")
	resp, err := ValidateToken(context.Background(), os.Getenv("FIREBASE_PROJECT_ID"), os.Getenv("TEST_JWT"))
	log.Println(resp, err)
}
