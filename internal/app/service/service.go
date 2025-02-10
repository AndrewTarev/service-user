package service

import (
	"context"

	"github.com/google/uuid"

	"service-user/internal/app/models"
	"service-user/internal/app/repository"
)

type ProfileService interface {
	CreateProfile(ctx context.Context, profile models.UserProfileInput) (uuid.UUID, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (models.UserProfileOut, error)
	UpdateProfile(ctx context.Context, profile models.UserProfileUpdate) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type Service struct {
	ProfileService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		ProfileService: NewServiceProfile(repo),
	}
}
