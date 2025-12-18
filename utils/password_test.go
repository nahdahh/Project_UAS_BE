package utils

import (
	"testing"
)

// TestHashPassword menguji fungsi HashPassword
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, 
		},
		{
			name:     "long password",
			password: "verylongpasswordwithmanycharsabcdefghijklmnopqrstuvwxyz123456789",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			hash, err := HashPassword(tt.password)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Error("HashPassword() returned empty hash")
			}
			if !tt.wantErr && hash == tt.password {
				t.Error("HashPassword() returned unhashed password")
			}
		})
	}
}

// TestCheckPasswordHash menguji fungsi CheckPasswordHash
func TestCheckPasswordHash(t *testing.T) {
	// ARRANGE: Create a known hash for "password123"
	password := "password123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal("failed to hash password for test setup")
	}

	tests := []struct {
		name     string
		hash     string
		password string
		want     bool
	}{
		{
			name:     "correct password",
			hash:     hash,
			password: "password123",
			want:     true,
		},
		{
			name:     "incorrect password",
			hash:     hash,
			password: "wrongpassword",
			want:     false,
		},
		{
			name:     "empty password",
			hash:     hash,
			password: "",
			want:     false,
		},
		{
			name:     "invalid hash",
			hash:     "invalid_hash",
			password: "password123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result := CheckPasswordHash(tt.hash, tt.password)

			// ASSERT
			if result != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestComparePassword menguji fungsi ComparePassword (alias)
func TestComparePassword(t *testing.T) {
	// ARRANGE
	password := "testpass"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal("failed to hash password for test setup")
	}

	tests := []struct {
		name     string
		hash     string
		password string
		want     bool
	}{
		{
			name:     "matching password",
			hash:     hash,
			password: "testpass",
			want:     true,
		},
		{
			name:     "non-matching password",
			hash:     hash,
			password: "wrongpass",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			result := ComparePassword(tt.hash, tt.password)

			// ASSERT
			if result != tt.want {
				t.Errorf("ComparePassword() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestHashPasswordConsistency menguji konsistensi hashing
func TestHashPasswordConsistency(t *testing.T) {
	// ARRANGE
	password := "consistencytest"

	// ACT: Hash same password twice
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	// ASSERT: Both should succeed
	if err1 != nil || err2 != nil {
		t.Fatal("HashPassword() returned error")
	}

	// ASSERT: Hashes should be different (bcrypt uses salt)
	if hash1 == hash2 {
		t.Error("HashPassword() returned identical hashes (should use different salts)")
	}

	// ASSERT: Both hashes should validate the password
	if !CheckPasswordHash(hash1, password) {
		t.Error("first hash doesn't validate password")
	}
	if !CheckPasswordHash(hash2, password) {
		t.Error("second hash doesn't validate password")
	}
}
