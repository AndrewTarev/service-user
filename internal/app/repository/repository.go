package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"service-user/internal/app/models"
)

// ProfileRepository - интерфейс репозитория для работы с профилем пользователя
type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile models.UserProfileInput) (uuid.UUID, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (models.UserProfileOut, error)
	UpdateProfile(ctx context.Context, profile models.UserProfileUpdate) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type Repository struct {
	ProfileRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		ProfileRepository: NewProfileRepository(db),
	}
}
