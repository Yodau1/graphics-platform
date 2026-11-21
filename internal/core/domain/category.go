package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	"github.com/google/uuid"
)

type Category struct {
	ID      uuid.UUID
	Version int64
	Name    string
}

func NewCategory(
	id uuid.UUID,
	version int64,
	name string,
) Category {
	return Category{
		ID:      id,
		Version: version,
		Name:    name,
	}
}

func CreateCategory(name string) Category {
	return NewCategory(uuid.New(), 1, strings.TrimSpace(name))
}

// Валидация на длину имени
func (c *Category) Validate() error {
	nameLength := utf8.RuneCountInString(c.Name)
	if nameLength < 1 || nameLength > 100 {
		return fmt.Errorf(
			"category name must contain 1 to 100 characters: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

// Содержит полне, переданное на изменение
type CategoryPatch struct {
	Name Nullable[string]
}

func NewCategoryPatch(name Nullable[string]) CategoryPatch {
	return CategoryPatch{
		Name: name,
	}
}

// Валидация, чтобы не было NULL
func (p *CategoryPatch) Validate() error {
	if !p.Name.Set {
		return fmt.Errorf("category patch must include name: %w", core_errors.ErrInvalidArgument)
	}

	if p.Name.Value == nil {
		return fmt.Errorf("name cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (c *Category) ApplyPatch(patch CategoryPatch) error {

	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate category patch: %w", err)
	}

	tmp := *c
	tmp.Name = strings.TrimSpace(*patch.Name.Value)

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched category: %w", err)
	}

	*c = tmp

	return nil
}
