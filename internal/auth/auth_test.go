package auth

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	// Test cases with different passwords
	testCases := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "C0mpl3x_P@ssw0rd!",
		},
		{
			name:     "empty password",
			password: "",
		},
		{
			name:     "very long password",
			password: "ThisIsAVeryLongPasswordThatShouldStillWork1234567890!@#$%^&*()",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(tc.password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			// Check that the hash is not empty
			if hash == "" {
				t.Error("HashPassword() returned an empty hash")
			}

			// Verify that the original password matches the hash
			err = CheckPasswordHash(tc.password, hash)
			if err != nil {
				t.Errorf("CheckPasswordHash() failed to verify a valid password: %v", err)
			}

			// Test with an incorrect password
			wrongPassword := tc.password + "wrong"
			err = CheckPasswordHash(wrongPassword, hash)
			if err == nil {
				t.Error("CheckPasswordHash() verified an incorrect password")
			}
		})
	}
}

func TestHashPassword_MultipleCallsProduceDifferentHashes(t *testing.T) {
	// Given the same password, multiple calls to HashPassword should produce different hashes
	// due to the salt that bcrypt adds
	password := "samepassword"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash1 == hash2 {
		t.Error("Multiple calls to HashPassword() with the same password produced identical hashes")
	}

	// Both hashes should still verify the same password
	if err := CheckPasswordHash(password, hash1); err != nil {
		t.Error("First hash failed to verify password")
	}

	if err := CheckPasswordHash(password, hash2); err != nil {
		t.Error("Second hash failed to verify password")
	}
}
