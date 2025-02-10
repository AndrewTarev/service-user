package service

import (
	"context"

	"github.com/google/uuid"

	"service-user/internal/app/models"
	"service-user/internal/app/repository"
)

type Profile struct {
	repo *repository.Repository
}

func NewServiceProfile(repo *repository.Repository) *Profile {
	return &Profile{repo}
}

func (p *Profile) CreateProfile(ctx context.Context, profile models.UserProfileInput) (uuid.UUID, error) {
	id, err := p.repo.CreateProfile(ctx, profile)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (p *Profile) GetProfile(ctx context.Context, userID uuid.UUID) (models.UserProfileOut, error) {
	profile, err := p.repo.GetProfile(ctx, userID)
	if err != nil {
		return models.UserProfileOut{}, err
	}
	return profile, nil
}

func (p *Profile) UpdateProfile(ctx context.Context, profile models.UserProfileUpdate) error {
	err := p.repo.UpdateProfile(ctx, profile)
	if err != nil {
		return err
	}
	return nil
}

func (p *Profile) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	err := p.repo.DeleteProfile(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
