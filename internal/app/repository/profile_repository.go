package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	logger "github.com/sirupsen/logrus"

	"service-user/internal/app/errs"
	"service-user/internal/app/models"
)

const (
	DuplicateValue = "23505"
)

// ProfileRepos - структура для работы с базой данных через pgxpool
type ProfileRepos struct {
	db *pgxpool.Pool
}

// NewProfileRepository - конструктор для создания нового репозитория
func NewProfileRepository(db *pgxpool.Pool) *ProfileRepos {
	return &ProfileRepos{db: db}
}

// CreateProfile - создание профиля пользователя
func (r *ProfileRepos) CreateProfile(ctx context.Context, profile models.UserProfileInput) (uuid.UUID, error) {
	query := `
		INSERT INTO user_profiles (user_id, first_name, last_name, city)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, profile.UserID, profile.FirstName, profile.LastName, profile.City).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == DuplicateValue {
				return uuid.UUID{}, errs.ErrProfileAlreadyExists
			}
			logger.Errorf("Error while inserting user-profile %v", err)
			return uuid.UUID{}, errs.ErrCreateUserProfile
		}
	}
	logger.Infof("Created user-profile %v", id)
	return id, nil
}

// GetProfile - получение профиля пользователя по userID
func (r *ProfileRepos) GetProfile(ctx context.Context, userID uuid.UUID) (models.UserProfileOut, error) {
	query := `SELECT id, user_id, first_name, last_name, city, created_at, updated_at FROM user_profiles WHERE user_id = $1`
	row := r.db.QueryRow(ctx, query, userID)

	var profile models.UserProfileOut
	err := row.Scan(&profile.ID, &profile.UserID, &profile.FirstName, &profile.LastName, &profile.City, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return models.UserProfileOut{}, errs.ErrProfileNotFound
		}
		logger.Errorf("Error while getting user-profile %v", err)
		return models.UserProfileOut{}, errs.ErrProfileNotFound
	}
	return profile, nil
}

// UpdateProfile - частичное обновление профиля пользователя
func (r *ProfileRepos) UpdateProfile(ctx context.Context, profile models.UserProfileUpdate) error {
	var updates []string
	var args []interface{}
	argID := 1

	// Проверяем, какие поля переданы, и добавляем их в SQL-запрос
	if profile.FirstName != "" {
		updates = append(updates, fmt.Sprintf("first_name = $%d", argID))
		args = append(args, profile.FirstName)
		argID++
	}
	if profile.LastName != "" {
		updates = append(updates, fmt.Sprintf("last_name = $%d", argID))
		args = append(args, profile.LastName)
		argID++
	}
	if profile.City != "" {
		updates = append(updates, fmt.Sprintf("city = $%d", argID))
		args = append(args, profile.City)
		argID++
	}

	// Добавляем user_id в конец списка аргументов
	args = append(args, profile.UserID)

	// Формируем SQL-запрос
	query := fmt.Sprintf(`
		UPDATE user_profiles 
		SET %s 
		WHERE user_id = $%d`, strings.Join(updates, ", "), argID)

	// Выполняем запрос
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		logger.Errorf("Error while updating user-profile %v", err)
		return errs.ErrUpdateUserProfile
	}

	return nil
}

// DeleteProfile - удаление профиля пользователя
func (r *ProfileRepos) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	// проверяем наличие профиля
	_, err := r.GetProfile(ctx, userID)
	if err != nil {
		return err
	}

	query := `DELETE FROM user_profiles WHERE user_id = $1`
	_, err = r.db.Exec(ctx, query, userID)
	if err != nil {
		logger.Errorf("error while deleting user-profile %v", err)
		return errs.ErrDeleteUserProfile
	}
	return nil
}
