package utils

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/CNMoreno/cnm-proyect-go/internal/constants"
	"github.com/CNMoreno/cnm-proyect-go/internal/domain"
	"github.com/gocarina/gocsv"
)

var OpenFileFunc = func(file *multipart.FileHeader) (multipart.File, error) {
	return file.Open()
}

// ReadCSVFile handles read csv and extract data and assing to user.
func ReadCSVFile(file *multipart.FileHeader) ([]domain.User, string) {
	extension := strings.ToLower(filepath.Ext(file.Filename))
	if extension != ".csv" {
		return nil, constants.ErrOnlyAcceptCSVFile
	}

	csvFile, err := OpenFileFunc(file)
	if err != nil {
		return nil, constants.ErrOpenFile
	}

	defer func() {
		err = csvFile.Close()
	}()

	if err != nil {
		return nil, constants.ErrClosingFile
	}

	var users []domain.User

	if err := gocsv.Unmarshal(csvFile, &users); err != nil {
		return nil, constants.ErrProcessCSVFile
	}

	return users, ""
}
