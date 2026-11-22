package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	"github.com/google/uuid"
)

//определяем поля карточки, её начальное состояние и редактируемые поля

type ResourceStatus string

const (
	ResourceStatusDraft     ResourceStatus = "draft"
	ResourceStatusPending   ResourceStatus = "pending"
	ResourceStatusPublished ResourceStatus = "published"
	ResourceStatusRejected  ResourceStatus = "rejected"
	ResourceStatusHidden    ResourceStatus = "hidden"
	ResourceStatusArchived  ResourceStatus = "archived"
)

type Resource struct {
	ID          uuid.UUID
	Version     int64
	AuthorID    uuid.UUID
	CategoryID  uuid.UUID
	Title       string
	Description string
	Status      ResourceStatus
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func NewResource(
	id uuid.UUID,
	version int64,
	authorID uuid.UUID,
	categoryID uuid.UUID,
	title string,
	description string,
	status ResourceStatus,
	createdAt time.Time,
	updatedAt *time.Time,
) Resource {
	return Resource{
		ID:          id,
		Version:     version,
		AuthorID:    authorID,
		CategoryID:  categoryID,
		Title:       title,
		Description: description,
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func CreateResource(
	authorID uuid.UUID,
	categoryID uuid.UUID,
	title string,
	description string,
) Resource {
	now := time.Now().UTC()

	return NewResource(
		uuid.New(),
		1,
		authorID,
		categoryID,
		strings.TrimSpace(title),
		description,
		ResourceStatusDraft,
		now,
		nil,
	)
}

func (r *Resource) Validate() error {
	if r.AuthorID == uuid.Nil {
		return fmt.Errorf("resource author is required: %w", core_errors.ErrInvalidArgument)
	}

	if r.CategoryID == uuid.Nil {
		return fmt.Errorf("resource category is required: %w", core_errors.ErrInvalidArgument)
	}

	titleLength := utf8.RuneCountInString(strings.TrimSpace(r.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf(
			"resource title must contain 1 to 100 characters: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type ResourcePatch struct {
	CategoryID  Nullable[uuid.UUID]
	Title       Nullable[string]
	Description Nullable[string]
}

func NewResourcePatch(
	categoryID Nullable[uuid.UUID],
	title Nullable[string],
	description Nullable[string],
) ResourcePatch {
	return ResourcePatch{
		CategoryID:  categoryID,
		Title:       title,
		Description: description,
	}
}

func (p *ResourcePatch) Validate() error {
	if !p.CategoryID.Set && !p.Title.Set && !p.Description.Set {
		return fmt.Errorf("resource patch must include a field: %w", core_errors.ErrInvalidArgument)
	}

	if p.CategoryID.Set && p.CategoryID.Value == nil {
		return fmt.Errorf("category_id cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}

	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("title cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}

	if p.Description.Set && p.Description.Value == nil {
		return fmt.Errorf("description cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (r *Resource) ApplyPatch(patch ResourcePatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate resource patch: %w", err)
	}

	if r.Status == ResourceStatusArchived {
		return fmt.Errorf("archived resource cannot be edited: %w", core_errors.ErrConflict)
	}

	tmp := *r

	if patch.CategoryID.Set {
		tmp.CategoryID = *patch.CategoryID.Value
	}

	if patch.Title.Set {
		tmp.Title = strings.TrimSpace(*patch.Title.Value)
	}

	if patch.Description.Set {
		tmp.Description = *patch.Description.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched resource: %w", err)
	}

	tmp.Status = ResourceStatusDraft
	updatedAt := time.Now().UTC()
	tmp.UpdatedAt = &updatedAt
	*r = tmp

	return nil
}

// Переводит карточку в архив
func (r *Resource) Archive() bool {
	if r.Status == ResourceStatusArchived {
		return false
	}

	r.Status = ResourceStatusArchived
	updatedAt := time.Now().UTC()
	r.UpdatedAt = &updatedAt

	return true
}
