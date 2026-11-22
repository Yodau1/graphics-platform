package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
	"github.com/google/uuid"
)

// Переводим пользователя в роль автора
func (r *UsersRepository) BecomeAuthor(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE graphics.users
		SET role = 'author',
			version = CASE WHEN role = 'user' THEN version + 1 ELSE version END
		WHERE id = $1 AND is_active AND role IN ('user', 'author')
		RETURNING id, name, email, role, is_active, version;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var user userModel
	if err := user.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("accout cannot become an author: %w", core_errors.ErrForbidden)
		}
		return domain.User{}, fmt.Errorf("scan updated user: %w", err)
	}

	userDomain := modelToDomain(user)
	return userDomain, nil
}
