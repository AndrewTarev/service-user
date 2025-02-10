package utils

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"

	logger "github.com/sirupsen/logrus"

	"service-user/internal/app/errs"
)

type JWTManager struct {
	publicKey *rsa.PublicKey
}

// NewJWTManager загружает RSA-ключи и создает JWT-менеджер
func NewJWTManager(publicKeyPath string) (*JWTManager, error) {
	// Загружаем публичный ключ
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		logger.WithError(err).Error("failed to read public key")
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		logger.WithError(err).Error("failed to parse public key")
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &JWTManager{
		publicKey: publicKey,
	}, nil
}

// DecodeJWT парсит токен и проверяет его подпись публичным ключом
func (j *JWTManager) DecodeJWT(tokenString string) (jwt.Claims, error) {
	// Разбираем и проверяем подпись токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется алгоритм RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil // Возвращаем публичный ключ для проверки подписи
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.Debugf("Error: token is expired: %v", err)
			return nil, errs.ErrTokenExpired
		}
		logger.Debugf("Error parsing token: %v", err)
		return nil, errs.ErrTokenInvalid
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errs.ErrTokenInvalid
}
