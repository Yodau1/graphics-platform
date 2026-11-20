package users_service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

const sessionTTL = 24 * time.Hour

// Возвращает токен сесии, время жизни и ошибку
func (s *UsersService) Login(
	ctx context.Context,
	email string, password string,
) (string, time.Time, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	//Пытаемся по email получить userID и хэш пароля зарегистрированного пользователя
	userID, passwordHash, isActive, err := s.usersRepository.GetUserCredentials(ctx, email)
	if errors.Is(err, core_errors.ErrNotFound) {
		return "", time.Time{}, core_errors.ErrUnauthorized
	}
	if err != nil {
		return "", time.Time{}, fmt.Errorf("get user credentials: %w", err)
	}

	//Сравниваем переданный пароль с хешем того, что лежит в БД
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordHash)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", time.Time{}, core_errors.ErrUnauthorized
		}
		return "", time.Time{}, fmt.Errorf("compare password hash: %w", err)
	}
	if !isActive {
		return "", time.Time{}, core_errors.ErrUnauthorized
	}

	//Если пользователь нашелся, верный пароль и его учетка активна, то выдаем токен сессии
	token := rand.Text()
	tokenHash := sha256.Sum256([]byte(token))
	//задаем срок жизни
	expiresAt := time.Now().UTC().Add(sessionTTL)

	if err := s.usersRepository.SaveUserSession(ctx, userID, tokenHash[:], expiresAt); err != nil {
		return "", time.Time{}, fmt.Errorf("save user session: %w", err)
	}

	return token, expiresAt, nil

}

// Предоставление доступа пользователю по токену сессии
func (s *UsersService) Authenticate(ctx context.Context, token string) (domain.User, error) {
	tokenHash := sha256.Sum256([]byte(token))
	user, err := s.usersRepository.GetUserBySession(ctx, tokenHash[:])
	if errors.Is(err, core_errors.ErrNotFound) {
		return domain.User{}, core_errors.ErrUnauthorized
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by session: %w", err)
	}

	return user, nil
}

// Отзываем токен, удаляя сессию
func (s *UsersService) Logout(ctx context.Context, token string) error {
	tokenHash := sha256.Sum256([]byte(token))

	if err := s.usersRepository.DeleteUserSession(ctx, tokenHash[:]); err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	return nil
}
