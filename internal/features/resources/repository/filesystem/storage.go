package resources_filesystem_repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

//Файловое хранилище
//Создаем каталог для ресурса

type ResourceFileStorage struct {
	//Хранит путь к общему каталогу
	root string
}

func NewResourceFileStorage(root string) *ResourceFileStorage {
	return &ResourceFileStorage{
		root: root,
	}
}

// Записывает содержимое под ключом, сгенерированным сервером
func (s *ResourceFileStorage) SaveFile(ctx context.Context, key string, data []byte) error {
	//Если контекст завершен, то вернуть ошибку
	if err := ctx.Err(); err != nil {
		return err
	}

	//Проверяет, можно ли прочитать строку key как uuid
	if _, err := uuid.Parse(key); err != nil {
		return fmt.Errorf("invalid storage key: %w", err)
	}

	if err := os.MkdirAll(s.root, 0700); err != nil {
		return fmt.Errorf("create resource files directory: %w", err)
	}

	path := filepath.Join(s.root, key)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("open resource file: %w", err)
	}

	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		//Если была ошибка, удаляем файл по пути
		return fmt.Errorf("write resource file: %w", errors.Join(err, os.Remove(path)))
	}

	return nil
}

// Используется для очистки файла при отклонённой записи
func (s *ResourceFileStorage) DeleteFile(key string) error {
	if _, err := uuid.Parse(key); err != nil {
		return fmt.Errorf("invalid storage key: %w", err)
	}

	err := os.Remove(filepath.Join(s.root, key))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}
