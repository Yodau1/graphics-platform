package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
)

// SaveUser вставляет нового пользователя в БД и возвращает сохранённую версию.
// RETURNING позволяет получить итоговое состояние записи одним запросом.
func (r *UsersRepository) SaveUser(
	ctx context.Context,
	user domain.User,
	passwordHash string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO graphics.users (id, name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO NOTHING
		RETURNING id, name, email, role, is_active, version;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		passwordHash,
	)

	var model userModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("email already registered: %w", core_errors.ErrConflict)
		}
		return domain.User{}, fmt.Errorf("scan saved user: %w", err)
	}
	userDomain := modelToDomain(model)
	return userDomain, nil
}
