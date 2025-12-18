package utils

import (
	"testing"
	"time"
)

// TestInitJWT menguji inisialisasi JWT secret
func TestInitJWT(t *testing.T) {
	// ARRANGE
	testSecret := "test_secret_key_12345"

	// ACT
	InitJWT(testSecret)

	// ASSERT
	if len(jwtSecret) == 0 {
		t.Error("InitJWT() failed to set jwtSecret")
	}
	if string(jwtSecret) != testSecret {
		t.Errorf("InitJWT() set wrong secret, got %s, want %s", string(jwtSecret), testSecret)
	}
}

// TestGenerateToken menguji pembuatan JWT token
func TestGenerateToken(t *testing.T) {
	// ARRANGE
	InitJWT("test_secret")

	tests := []struct {
		name        string
		userID      string
		username    string
		email       string
		role        string
		permissions []string
		duration    time.Duration
		wantErr     bool
	}{
		{
			name:        "valid token generation",
			userID:      "user123",
			username:    "testuser",
			email:       "test@example.com",
			role:        "Admin",
			permissions: []string{"read", "write"},
			duration:    24 * time.Hour,
			wantErr:     false,
		},
		{
			name:        "token with empty permissions",
			userID:      "user456",
			username:    "student",
			email:       "student@example.com",
			role:        "Mahasiswa",
			permissions: []string{},
			duration:    1 * time.Hour,
			wantErr:     false,
		},
		{
			name:        "token with short duration",
			userID:      "user789",
			username:    "lecturer",
			email:       "lecturer@example.com",
			role:        "Dosen Wali",
			permissions: []string{"verify_achievement"},
			duration:    5 * time.Minute,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			token, err := GenerateToken(tt.userID, tt.username, tt.email, tt.role, tt.permissions, tt.duration)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() returned empty token")
			}
		})
	}
}

// TestGenerateTokenWithoutInit menguji token generation tanpa inisialisasi
func TestGenerateTokenWithoutInit(t *testing.T) {
	// ARRANGE: Reset jwtSecret
	jwtSecret = []byte{}

	// ACT
	token, err := GenerateToken("user1", "test", "test@test.com", "Admin", []string{}, 1*time.Hour)

	// ASSERT
	if err == nil {
		t.Error("GenerateToken() should return error when JWT not initialized")
	}
	if token != "" {
		t.Error("GenerateToken() should return empty token on error")
	}

	// Cleanup: Re-initialize for other tests
	InitJWT("test_secret")
}

// TestVerifyToken menguji verifikasi JWT token
func TestVerifyToken(t *testing.T) {
	// ARRANGE
	InitJWT("test_secret_verify")
	userID := "user123"
	username := "testuser"
	email := "test@example.com"
	role := "Admin"
	permissions := []string{"read", "write"}

	validToken, err := GenerateToken(userID, username, email, role, permissions, 1*time.Hour)
	if err != nil {
		t.Fatal("failed to generate token for test setup")
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "invalid token format",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			claims, err := VerifyToken(tt.token)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if claims == nil {
					t.Error("VerifyToken() returned nil claims for valid token")
				}
				if claims["sub"] != userID {
					t.Errorf("VerifyToken() sub = %v, want %v", claims["sub"], userID)
				}
				if claims["role"] != role {
					t.Errorf("VerifyToken() role = %v, want %v", claims["role"], role)
				}
			}
		})
	}
}

// TestGetClaimsFromToken menguji ekstraksi claims dari token
func TestGetClaimsFromToken(t *testing.T) {
	// ARRANGE
	InitJWT("test_secret_claims")
	userID := "user456"
	username := "claimuser"
	email := "claim@example.com"
	role := "Mahasiswa"
	permissions := []string{"create_achievement"}

	validToken, err := GenerateToken(userID, username, email, role, permissions, 2*time.Hour)
	if err != nil {
		t.Fatal("failed to generate token for test setup")
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token with claims",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			claims, err := GetClaimsFromToken(tt.token)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClaimsFromToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if claims == nil {
					t.Fatal("GetClaimsFromToken() returned nil claims")
				}
				if claims.Sub != userID {
					t.Errorf("GetClaimsFromToken() Sub = %v, want %v", claims.Sub, userID)
				}
				if claims.Username != username {
					t.Errorf("GetClaimsFromToken() Username = %v, want %v", claims.Username, username)
				}
				if claims.Email != email {
					t.Errorf("GetClaimsFromToken() Email = %v, want %v", claims.Email, email)
				}
				if claims.Role != role {
					t.Errorf("GetClaimsFromToken() Role = %v, want %v", claims.Role, role)
				}
				if len(claims.Permissions) != len(permissions) {
					t.Errorf("GetClaimsFromToken() Permissions count = %v, want %v", len(claims.Permissions), len(permissions))
				}
			}
		})
	}
}

// TestTokenExpiration menguji token yang sudah expired
func TestTokenExpiration(t *testing.T) {
	// ARRANGE
	InitJWT("test_secret_expiration")

	// Create token with very short duration
	expiredToken, err := GenerateToken("user1", "test", "test@test.com", "Admin", []string{}, 1*time.Millisecond)
	if err != nil {
		t.Fatal("failed to generate token for test setup")
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// ACT
	_, err = VerifyToken(expiredToken)

	// ASSERT
	if err == nil {
		t.Error("VerifyToken() should return error for expired token")
	}
}

// TestVerifyTokenWithoutInit menguji verifikasi tanpa inisialisasi
func TestVerifyTokenWithoutInit(t *testing.T) {
	// ARRANGE: Reset jwtSecret
	jwtSecret = []byte{}

	// ACT
	claims, err := VerifyToken("some.token.here")

	// ASSERT
	if err == nil {
		t.Error("VerifyToken() should return error when JWT not initialized")
	}
	if claims != nil {
		t.Error("VerifyToken() should return nil claims on error")
	}

	// Cleanup
	InitJWT("test_secret")
}

// TestTokenRoundTrip menguji generate dan verify token secara lengkap
func TestTokenRoundTrip(t *testing.T) {
	// ARRANGE
	InitJWT("roundtrip_secret")
	userID := "roundtrip_user"
	username := "roundtripuser"
	email := "roundtrip@test.com"
	role := "Dosen Wali"
	permissions := []string{"verify", "reject"}

	// ACT: Generate token
	token, err := GenerateToken(userID, username, email, role, permissions, 24*time.Hour)
	if err != nil {
		t.Fatal("GenerateToken() failed:", err)
	}

	// ACT: Verify token
	claims, err := VerifyToken(token)
	if err != nil {
		t.Fatal("VerifyToken() failed:", err)
	}

	// ASSERT: Check all claims
	if claims["sub"] != userID {
		t.Errorf("sub claim = %v, want %v", claims["sub"], userID)
	}
	if claims["username"] != username {
		t.Errorf("username claim = %v, want %v", claims["username"], username)
	}
	if claims["email"] != email {
		t.Errorf("email claim = %v, want %v", claims["email"], email)
	}
	if claims["role"] != role {
		t.Errorf("role claim = %v, want %v", claims["role"], role)
	}

	// Check permissions (it's stored as interface{})
	permsInterface := claims["permissions"]
	if permsInterface == nil {
		t.Error("permissions claim is nil")
	}
}

// TestGetClaimsFromExpiredToken menguji ekstraksi claims dari expired token
func TestGetClaimsFromExpiredToken(t *testing.T) {
	// ARRANGE
	InitJWT("test_expired_claims")

	// Create expired token
	expiredToken, err := GenerateToken("user1", "test", "test@test.com", "Admin", []string{}, 1*time.Nanosecond)
	if err != nil {
		t.Fatal("failed to generate token")
	}

	time.Sleep(10 * time.Millisecond)

	// ACT
	claims, err := GetClaimsFromToken(expiredToken)

	// ASSERT
	if err == nil {
		t.Error("GetClaimsFromToken() should return error for expired token")
	}
	if claims != nil {
		t.Error("GetClaimsFromToken() should return nil for expired token")
	}
}
