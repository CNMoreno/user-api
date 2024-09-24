package repository

import "golang.org/x/crypto/bcrypt"

// AppCrypto interface defines the methods for password hashing and comparison.
type AppCrypto interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

// BcryptCrypto struct implements the AppCrypto interface using bcrypt.
type BcryptCrypto struct{}

// GenerateFromPassword wraps bcrypt's GenerateFromPassword function.
func (BcryptCrypto) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// CompareHashAndPassword wraps bcrypt's CompareHashAndPassword function.
func (BcryptCrypto) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
