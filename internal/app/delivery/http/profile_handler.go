package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"service-user/internal/app/errs"
	"service-user/internal/app/models"
	"service-user/internal/app/service"
)

type ProfileHandler struct {
	service service.Service
}

func NewProfileHandler(service service.Service) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

// @securityDefinitions.cookie CookieAuth
// @name access_token
// @in cookie

// CreateProfile - создание профиля пользователя
// @Summary Создает профиль пользователя
// @Description Создает новый профиль пользователя, если он не существует
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param input body models.UserProfileInput true "Информация о профиле"
// @Security CookieAuth
// @Success 201 {object} models.ProfileIdResponse "Успешно создан профиль"
// @Failure 400 {object} middleware.ValidationErrorResponse "Неверный запрос"
// @Failure 401 {object} middleware.ValidationErrorResponse "Не авторизован"
// @Failure 500 {object} middleware.ValidationErrorResponse "Внутренняя ошибка сервера"
// @Router /user-profile/ [post]
func (ph *ProfileHandler) CreateProfile(c *gin.Context) {
	var input models.UserProfileInput

	// Парсим JSON в структуру
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	// Валидация входных данных
	if err := input.Validate(); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			c.Error(validationErrs)
			return
		}
	}

	userID := ph.GetUserIdFromContext(c)

	input.UserID = userID

	// Вызываем сервис создания профиля
	createdUserID, err := ph.service.CreateProfile(c, input)
	if err != nil {
		c.Error(err) // Передаем ошибку в middleware
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusCreated, models.ProfileIdResponse{
		ID: createdUserID,
	})
}

// @securityDefinitions.cookie CookieAuth
// @name access_token
// @in cookie

// GetProfile - получение профиля пользователя
// @Summary Получить профиль пользователя
// @Description Получает профиль пользователя по его ID
// @Tags Profile
// @Produce  json
// @Security CookieAuth
// @Success 200 {object} models.UserProfileOut "Информация о профиле"
// @Failure 401 {object} middleware.ValidationErrorResponse "Не авторизован"
// @Failure 404 {object} middleware.ValidationErrorResponse "Профиль не найден"
// @Failure 500 {object} middleware.ValidationErrorResponse "Внутренняя ошибка сервера"
// @Router /user-profile/ [get]
func (ph *ProfileHandler) GetProfile(c *gin.Context) {
	userID := ph.GetUserIdFromContext(c)
	userProfile, err := ph.service.GetProfile(c, userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, userProfile)
}

func (ph *ProfileHandler) GetUserIdFromContext(c *gin.Context) uuid.UUID {
	// Получаем user_id из контекста
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.Error(errs.ErrUnauthorized)
		return uuid.UUID{}
	}

	// Приводим userID к строке и парсим UUID
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.Error(errs.ErrInvalidUserId)
		return uuid.UUID{}
	}
	return userID
}

// @securityDefinitions.cookie CookieAuth
// @name access_token
// @in cookie

// UpdateProfile - обновление профиля пользователя
// @Summary Обновить профиль пользователя
// @Description Обновляет данные профиля пользователя
// @Tags Profile
// @Accept  json
// @Produce  json
// @Security CookieAuth
// @Param input body models.UserProfileUpdate true "Новые данные профиля"
// @Success 200 {object} models.SuccessResponse "Профиль успешно обновлен"
// @Failure 400 {object} middleware.ValidationErrorResponse "Неверный запрос"
// @Failure 401 {object} middleware.ValidationErrorResponse "Не авторизован"
// @Failure 500 {object} middleware.ValidationErrorResponse "Внутренняя ошибка сервера"
// @Router /user-profile/ [patch]
func (ph *ProfileHandler) UpdateProfile(c *gin.Context) {
	var input models.UserProfileUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	if err := input.Validate(); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			c.Error(validationErrs)
			return
		}
	}

	userID := ph.GetUserIdFromContext(c)
	input.UserID = userID

	err := ph.service.UpdateProfile(c, input)
	if err != nil {
		c.Error(err)
		return
	}

	response := models.SuccessResponse{
		Status: http.StatusOK,
		Data:   "Profile updated successfully",
	}

	c.JSON(http.StatusOK, response)
}

// @securityDefinitions.cookie CookieAuth
// @name access_token
// @in cookie

// DeleteProfile - удаление профиля пользователя
// @Summary Удалить профиль пользователя
// @Description Удаляет профиль пользователя по его ID
// @Tags Profile
// @Security CookieAuth
// @Success 204 {object} models.SuccessResponse"Профиль успешно удален"
// @Failure 401 {object} middleware.ValidationErrorResponse "Не авторизован"
// @Failure 404 {object} middleware.ValidationErrorResponse "Ошибка при удалении пользователя"
// @Failure 500 {object} middleware.ValidationErrorResponse "Внутренняя ошибка сервера"
// @Router /user-profile/ [delete]
func (ph *ProfileHandler) DeleteProfile(c *gin.Context) {
	userID := ph.GetUserIdFromContext(c)
	err := ph.service.DeleteProfile(c, userID)
	if err != nil {
		c.Error(err)
		return
	}
	response := models.SuccessResponse{
		Status: http.StatusNoContent,
		Data:   "Profile deleted successfully",
	}
	c.JSON(http.StatusOK, response)
}
