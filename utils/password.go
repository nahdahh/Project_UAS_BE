package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengubah password plain text menjadi hash menggunakan bcrypt
// password: password yang ingin di-hash
// return: hashed password atau error jika gagal
func HashPassword(password string) (string, error) {
	// GenerateFromPassword hash password dengan default cost (10)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash membandingkan password hash dengan password plain text
// hash: password yang sudah di-hash
// password: password plain text untuk dibandingkan
// return: true jika cocok, false jika tidak
func CheckPasswordHash(hash, password string) bool {
	// CompareHashAndPassword return error jika tidak cocok
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ComparePassword adalah alias dari CheckPasswordHash untuk compatibility
func ComparePassword(hash, password string) bool {
	return CheckPasswordHash(hash, password)
}

// GenerateHash adalah helper function untuk generate bcrypt hash
// Bisa dipanggil dari main.go atau script untuk testing
func GenerateHash(password string) {
	hash, err := HashPassword(password)
	if err != nil {
		panic(err)
	}
	println("Password:", password)
	println("Hash:", hash)
}
