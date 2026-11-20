package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
	"github.com/google/uuid"
)

// GetUserCredentials возвращает ID активного пользователя и хеш его пароля по email
func (r *UsersRepository) GetUserCredentials(
	ctx context.Context,
	email string,
) (uuid.UUID, string, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, password_hash, is_active
		FROM graphics.users
		WHERE email=$1;
	`

	var (
		id           uuid.UUID
		passwordHash string
		isActive     bool
	)

	row := r.pool.QueryRow(ctx, query, email)
	if err := row.Scan(&id, &passwordHash, &isActive); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return uuid.Nil, "", false, core_errors.ErrNotFound
		}
		return uuid.Nil, "", false, fmt.Errorf("scan user credentials: %w", err)
	}

	return id, passwordHash, isActive, nil
}

func (r *UsersRepository) SaveUserSession(
	ctx context.Context,
	userID uuid.UUID,
	tokenHash []byte,
	//время окончания сессии
	expiresAt time.Time,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	//Сохраняем хэш сессии
	query := `
		INSERT INTO graphics.user_sessions (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3);
	`

	if _, err := r.pool.Exec(ctx, query, tokenHash, userID, expiresAt); err != nil {
		return fmt.Errorf("insert user session: %w", err)
	}

	return nil
}

// GetUserBySession возвращает активного пользователя по хешу токена действующей сессии.
func (r *UsersRepository) GetUserBySession(
	ctx context.Context,
	tokenHash []byte,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	// Проверяем срок сессии и текущее состояние аккаунта
	query := `
		SELECT u.id, u.name, u.email, u.role, u.is_active, u.version
		FROM graphics.user_sessions AS s
		JOIN graphics.users AS u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > now() AND u.is_active;
	`

	var model userModel
	if err := model.Scan(r.pool.QueryRow(ctx, query, tokenHash)); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, core_errors.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("scan user by session: %w", err)
	}

	return modelToDomain(model), nil
}

// DeleteUserSession удаляет сессию по хешу её токена.
func (r *UsersRepository) DeleteUserSession(ctx context.Context, tokenHash []byte) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM graphics.user_sessions
		WHERE token_hash = $1;
	`

	if _, err := r.pool.Exec(ctx, query, tokenHash); err != nil {
		return fmt.Errorf("delete user session: %w", err)
	}

	return nil
}
