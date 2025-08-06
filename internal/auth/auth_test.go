package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "test-password"

	// Test password hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if hash == password {
		t.Fatal("Hash should not be the same as the original password")
	}

	// Test correct password check
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Failed to validate correct password: %v", err)
	}

	// Test incorrect password check
	err = CheckPasswordHash("wrong-password", hash)
	if err == nil {
		t.Fatal("Should have failed with incorrect password")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := time.Hour

	// Test creating JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	if token == "" {
		t.Fatal("Token should not be empty")
	}

	// Test validating JWT
	extractedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	if extractedUserID != userID {
		t.Fatalf("Extracted user ID (%s) doesn't match original user ID (%s)", extractedUserID, userID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"

	// Create token that expires immediately
	expiresIn := -1 * time.Hour // Expired 1 hour ago
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create expired JWT: %v", err)
	}

	// Test validating expired JWT
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("Should have failed with expired token")
	}

	// Check that the error is specifically about token expiration
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatalf("Expected TokenExpiredError but got: %v", err)
	}

}

func TestInvalidSecretJWT(t *testing.T) {
	userID := uuid.New()
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"
	expiresIn := time.Hour

	// Create token with correct secret
	token, err := MakeJWT(userID, correctSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Test validating with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Should have failed with wrong secret")
	}
}

func TestMalformedJWT(t *testing.T) {
	tokenSecret := "test-secret"
	malformedToken := "not.a.validtoken"

	// Test validating malformed JWT
	_, err := ValidateJWT(malformedToken, tokenSecret)
	if err == nil {
		t.Fatal("Should have failed with malformed token")
	}
}

func TestTokenWithInvalidUserID(t *testing.T) {
	// Create a token manually with an invalid UUID in the subject claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "Chirpy",
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
		"sub": "not-a-valid-uuid",
	})

	tokenSecret := "test-secret"
	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Validate the token - should fail parsing the UUID
	_, err = ValidateJWT(signedString, tokenSecret)
	if err == nil {
		t.Fatal("Should have failed with invalid UUID")
	}
}

func TestTokenWithMissingClaims(t *testing.T) {
	// Create a token without the subject claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "Chirpy",
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
		// Missing "sub" claim
	})

	tokenSecret := "test-secret"
	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Validate the token - should fail because subject claim is missing
	_, err = ValidateJWT(signedString, tokenSecret)
	if err == nil {
		t.Fatal("Should have failed with missing subject claim")
	}
}
