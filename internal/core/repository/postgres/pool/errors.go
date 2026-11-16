package core_postgres_pool

import "errors"

// Ошибки получаемые при выполнении запроса, которые позже репозиторий переведет в свои ошибки
var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown")
)
