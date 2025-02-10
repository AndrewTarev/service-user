package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// SuccessResponse - схема успешного ответа
type SuccessResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// UserProfileInput представляет профиль пользователя
type UserProfileInput struct {
	UserID    uuid.UUID `json:"-"` // Берется из JWT
	FirstName string    `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"last_name" validate:"required,min=2,max=50"`
	City      string    `json:"city" validate:"required,min=2,max=100"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (u *UserProfileInput) Validate() error {
	return validate.Struct(u)
}

type UserProfileUpdate struct {
	UserID    uuid.UUID `json:"-" validate:"required"` // Берется из JWT
	FirstName string    `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string    `json:"last_name" validate:"omitempty,min=2,max=50"`
	City      string    `json:"city" validate:"omitempty,min=2,max=100"`
}

func (u *UserProfileUpdate) Validate() error {
	return validate.Struct(u)
}

type UserProfileOut struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProfileIdResponse struct {
	ID uuid.UUID `json:"id"`
}

// Cart представляет корзину пользователя
type Cart struct {
	ID            uuid.UUID `json:"id"`
	UserProfileID uuid.UUID `json:"user_profile_id" validate:"required"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CartItem представляет элемент корзины
type CartItem struct {
	ID        uuid.UUID `json:"id"`
	CartID    uuid.UUID `json:"cart_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	Price     float64   `json:"price" validate:"required,gt=0"`
	CreatedAt time.Time `json:"created_at"`
}
