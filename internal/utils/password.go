package utils

import (
	"log"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generate a hash for password with bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compare plane password with his hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

// RegisterCustomValidators handles to register validation.
func RegisterCustomValidators(validate *validator.Validate) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("password", passwordValidator); err != nil {
			log.Fatalf("Error registering password validation %v", err)
		}
	}
}

func passwordValidator(f1 validator.FieldLevel) bool {
	password := f1.Field().String()

	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecialCharacter := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)

	return hasLowerCase && hasUpperCase && hasNumber && hasSpecialCharacter
}
