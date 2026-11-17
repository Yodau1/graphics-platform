package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

// requestValidator — глобальный экземпляр валидатора (go-playground/validator).
// Используется для валидации struct-тегов вида `validate:"required,min=1,max=100"`.
// Создаётся один раз на пакет (потокобезопасен).
var requestValidator = validator.New()

// validatable — интерфейс для DTO с собственным методом Validate().
// Если тип реализует этот интерфейс, используется его метод вместо struct-тегов.
// Это нужно для сложных правил валидации, которые нельзя выразить тегами
// (например, Nullable-поля в PATCH-запросах).
type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var err error
	if value, ok := dest.(validatable); ok {
		err = value.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
