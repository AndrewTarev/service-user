package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validateBankCard *validator.Validate

func init() {
	validate = validator.New()
}

// UserBankCard представляет банковскую карту пользователя
type UserBankCard struct {
	ID             uuid.UUID `json:"id"`
	UserProfileID  uuid.UUID `json:"-" validate:"required"`
	CardNumber     string    `json:"card_number" validate:"required,len=16,numeric"`
	ExpirationDate string    `json:"expiration_date" validate:"required,len=5,datetime=02/06"` // MM/YY
	CardHolderName string    `json:"card_holder_name" validate:"required,min=2,max=100"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u *UserBankCard) Validate() error {
	return validateBankCard.Struct(u)
}
