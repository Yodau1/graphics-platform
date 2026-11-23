package domain

import (
	"time"

	"github.com/google/uuid"
)

type ResourceFileKind string

const (
	//Сам файл для скачивания
	ResourceFileKindOriginal ResourceFileKind = "original"
	//Изображение для карточки
	ResourceFileKindImage ResourceFileKind = "image"
	//3D превью для карточки
	ResourceFileKindPreview3D ResourceFileKind = "preview_3d"
)

//У карточки «3D-модель стула» могут быть три файла:
//chair.glb — original, полная модель для скачивания
//chair.jpg — image, обложка для каталога
//chair-lite.glb — preview_3d, облегчённая модель, которую можно покрутить в браузере

// Структура описывает тип ресурса, который относится к карточке
type ResourceFile struct {
	ID uuid.UUID
	//ID карточки, к которой прикреплен файл
	ResourceID   uuid.UUID
	Kind         ResourceFileKind
	StorageKey   string
	OriginalName string
	Format       string
	SizeBytes    int64
	CreatedAt    time.Time
}

func CreateResourceFile(
	resourceID uuid.UUID,
	kind ResourceFileKind,
	originalName string,
	format string,
	sizeBytes int64,
) ResourceFile {
	id := uuid.New()
	return ResourceFile{
		ID:           id,
		ResourceID:   resourceID,
		Kind:         kind,
		StorageKey:   id.String(),
		OriginalName: originalName,
		Format:       format,
		SizeBytes:    sizeBytes,
		CreatedAt:    time.Now().UTC(),
	}
}
