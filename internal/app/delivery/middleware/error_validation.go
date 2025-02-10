package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	logger "github.com/sirupsen/logrus"

	"service-user/internal/app/errs"
)

// ValidationErrorResponse структура для JSON-ответа
type ValidationErrorResponse struct {
	Error struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Fields  map[string]string `json:"fields,omitempty"` // Поля с ошибками
	} `json:"error"`
}

// ErrorHandler глобальный middleware для обработки ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var statusCode int
			var message string
			var fieldErrors map[string]string // Словарь ошибок валидации

			var validationErrs validator.ValidationErrors // Объявляем переменную перед switch

			switch {
			case errors.Is(err, errs.ErrUnauthorized):
				statusCode = http.StatusUnauthorized
				message = "Unauthorized"
			case errors.Is(err, errs.ErrInvalidUserId):
				statusCode = http.StatusBadRequest
				message = "Invalid UserId"
			case errors.Is(err, errs.ErrProfileNotFound):
				statusCode = http.StatusNotFound
				message = "Profile not found"
			case errors.Is(err, errs.ErrDeleteUserProfile):
				statusCode = http.StatusNotFound
				message = "Error delete User Profile"
			case errors.Is(err, errs.ErrUpdateUserProfile):
				statusCode = http.StatusBadRequest
				message = "Error update User Profile"
			case errors.Is(err, errs.ErrTokenExpired):
				statusCode = http.StatusUnauthorized
				message = "Token is expired"
			case errors.Is(err, errs.ErrCreateUserProfile):
				statusCode = http.StatusBadRequest
				message = "User profile is invalid"
			case errors.Is(err, errs.ErrTokenInvalid):
				statusCode = http.StatusUnauthorized
				message = "Token is invalid"
			case errors.Is(err, errs.ErrProfileAlreadyExists):
				statusCode = http.StatusBadRequest
				message = "Profile already exists"
			case errors.As(err, &validationErrs): // Проверяем, является ли err ошибкой валидации
				statusCode = http.StatusBadRequest
				message = "Validation error"
				fieldErrors = make(map[string]string)
				for _, fieldErr := range validationErrs {
					fieldErrors[fieldErr.Field()] = validationErrorMessage(fieldErr)
				}

			default:
				statusCode = http.StatusInternalServerError
				message = "Internal server error"
			}

			// Логируем критические ошибки
			if statusCode == http.StatusInternalServerError {
				logger.Errorf("Unhandled server error: %v", err)
			}

			// Формируем JSON-ответ
			errorResponse := ValidationErrorResponse{}
			errorResponse.Error.Code = statusCode
			errorResponse.Error.Message = message
			if len(fieldErrors) > 0 {
				errorResponse.Error.Fields = fieldErrors
			}

			// Отправляем JSON-ответ
			c.JSON(statusCode, errorResponse)
		}
	}
}

// validationErrorMessage формирует читаемое сообщение ошибки
func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min":
		return "must be at least " + fe.Param()
	case "max":
		return "must be at most " + fe.Param()
	case "gt":
		return "must be greater than " + fe.Param()
	default:
		return "is invalid"
	}
}
