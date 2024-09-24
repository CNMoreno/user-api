package utils_test

import (
	"errors"
	"testing"

	"github.com/CNMoreno/cnm-proyect-go/internal/utils"
	mocks "github.com/CNMoreno/cnm-proyect-go/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type valuesTestCases struct {
	name         string
	password     string
	hashPassword string
	isError      bool
	err          error
}

func TestPassword(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name:         "should return hash password",
			password:     "test",
			hashPassword: "hashPassword",
		},
		{
			name:     "should throw an error when hashPasswords fails",
			password: "test",
			err:      errors.New("hashPassword Error"),
			isError:  true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockPassword := new(mocks.AppCrypto)

			mockPassword.On("GenerateFromPassword", mock.Anything, mock.Anything).Return([]byte(test.hashPassword), test.err)

			passwordUtils := utils.NewHashPassword(mockPassword)

			password, err := passwordUtils.HashPassword(test.password)

			if test.isError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.NotEmpty(t, password)
				assert.Equal(t, password, test.hashPassword)
			}

		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	testCases := []valuesTestCases{
		{
			name:         "should return true when check passoword is sucessful",
			password:     "test",
			hashPassword: "hashPassword",
		},
		{
			name:     "should return false when password is different",
			password: "test",
			err:      errors.New("check password Error"),
			isError:  true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockPassword := new(mocks.AppCrypto)

			mockPassword.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(test.err)

			passwordUtils := utils.NewHashPassword(mockPassword)

			result := passwordUtils.CheckPasswordHash(test.password, test.hashPassword)

			if test.isError {
				assert.False(t, result)
			} else {
				assert.True(t, result)
			}

		})
	}
}
