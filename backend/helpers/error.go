package helpers

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// TranslateErrorMessage mengubah error (validasi / database) menjadi map field -> pesan
func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)

	// Error validasi dari binding Gin
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorsMap[field] = field + " is required"
			case "email":
				errorsMap[field] = "Invalid email format"
			case "min":
				errorsMap[field] = field + " must be at least " + fieldError.Param() + " characters"
			case "max":
				errorsMap[field] = field + " must be at most " + fieldError.Param() + " characters"
			default:
				errorsMap[field] = "Invalid value for " + field
			}
		}
		return errorsMap
	}

	// Error dari database
	if err != nil {
		switch {
		case IsDuplicateEntryError(err):
			if strings.Contains(err.Error(), "username") {
				errorsMap["Username"] = "Username already exists"
			} else if strings.Contains(err.Error(), "email") {
				errorsMap["Email"] = "Email already exists"
			} else {
				errorsMap["Error"] = "Duplicate entry"
			}
		case errors.Is(err, gorm.ErrRecordNotFound):
			errorsMap["Error"] = "Record not found"
		default:
			errorsMap["Error"] = err.Error()
		}
	}

	return errorsMap
}

// IsDuplicateEntryError mengecek apakah error berasal dari duplicate key (MySQL 1062)
func IsDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(err.Error(), "Duplicate entry")
}
