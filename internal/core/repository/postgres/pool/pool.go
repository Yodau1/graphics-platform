package core_postgres_pool

import (
	"context"
	"time"
)

type Pool interface {
	//Читает несколько строк и возвращает Rows
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	//Выполняет запроса, возвращает Row
	QueryRow(ctx context.Context, sql string, args ...any) Row
	//Выполнение команды без чтения результата
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Close()

	OpTimeout() time.Duration
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type Row interface {
	Scan(dest ...any) error
}

// CommandTag — результат запроса
type CommandTag interface {
	//Кол-во удаленных строк
	RowsAffected() int64
}
