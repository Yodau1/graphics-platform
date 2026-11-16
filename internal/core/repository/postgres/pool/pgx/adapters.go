package core_pgx_pool

import (
	"errors"
	"fmt"

	core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// pgxRows оборачивает pgx.Rows для реализации интерфейса core_postgres_pool.Rows
type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}
	return nil
}

// pgxCommandTag оборачивает pgconn.CommandTag для реализации интерфейса CommandTag.
type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	//Ошбика PostgreSQL "нарушение внешнего ключа"
	const pgxViolatesForeignKeyErrorCode = "23503"

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	// Проверяем, является ли ошибка структурированной PostgreSQL-ошибкой.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgxViolatesForeignKeyErrorCode {
			return fmt.Errorf(
				"%w: %w",
				err,
				core_postgres_pool.ErrViolatesForeignKey,
			)
		}
	}

	// Все остальные ошибки оборачиваем в ErrUnknown.
	return fmt.Errorf(
		"%v: %w",
		err,
		core_postgres_pool.ErrUnknown,
	)
}
