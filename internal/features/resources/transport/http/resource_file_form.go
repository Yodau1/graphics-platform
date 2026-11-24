package resources_transport_http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Yodau1/graphics-platform/internal/core/domain"
	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
	resources_service "github.com/Yodau1/graphics-platform/internal/features/resources/service"
)

type resourceFileForm struct {
	kind         domain.ResourceFileKind
	originalName string
	data         []byte
}

//Этот файл отвечает за формат HTTP-запроса:
// лимит тела, чтение файла и удаление временных файлов multipart

func readResourceFileForm(
	rw http.ResponseWriter,
	r *http.Request,
	requireKind bool,
) (resourceFileForm, error) {

	r.Body = http.MaxBytesReader(rw, r.Body, resources_service.MaxResourceFileSize+64*1024)

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		return resourceFileForm{}, fmt.Errorf(
			"invalid multipart form or file too large: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	defer r.MultipartForm.RemoveAll() // Удаляем временные файлы формы при выходе.

	// len проверяет, что в поле file передан ровно один файл.
	files := r.MultipartForm.File["file"]
	if len(files) != 1 {
		return resourceFileForm{}, fmt.Errorf("one file is required: %w", core_errors.ErrInvalidArgument)
	}

	header := files[0]
	if header.Size > resources_service.MaxResourceFileSize {
		return resourceFileForm{}, fmt.Errorf("file exceeds 10 MiB: %w", core_errors.ErrInvalidArgument)
	}

	form := resourceFileForm{originalName: header.Filename}
	if requireKind {
		kinds := r.MultipartForm.Value["kind"]
		if len(kinds) != 1 {
			return resourceFileForm{}, fmt.Errorf("one kind is required: %w", core_errors.ErrInvalidArgument)
		}

		form.kind = domain.ResourceFileKind(kinds[0])
	}

	source, err := header.Open()
	if err != nil {
		return resourceFileForm{}, fmt.Errorf("open uploaded file: %w", err)
	}
	defer source.Close()

	form.data, err = io.ReadAll(io.LimitReader(source, resources_service.MaxResourceFileSize+1))
	if err != nil {
		return resourceFileForm{}, fmt.Errorf("read uploaded file: %w", err)
	}

	return form, nil
}
