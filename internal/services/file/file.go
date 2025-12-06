package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"main/internal/entity"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Upload(ctx context.Context, image *multipart.FileHeader, folder string) (entity.File, error) {
	if image == nil {
		return entity.File{}, errors.New("image not found")
	}

	mimeType := image.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		return entity.File{}, fmt.Errorf("invalid file type: expected image, got %s", mimeType)
	}

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return entity.File{}, err
	}

	randomName := uuid.New().String() + filepath.Ext(image.Filename)
	dst := filepath.Join(folder, randomName)

	src, err := image.Open()
	if err != nil {
		return entity.File{}, err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return entity.File{}, err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return entity.File{}, err
	}

	info := entity.File{
		OriginalName: image.Filename,
		SavedName:    randomName,
		Size:         fmt.Sprintf("%d", image.Size),
		Path:         dst,
	}

	return info, nil
}

func (s Service) Delete(ctx context.Context, url string) error {
	err := os.Remove("./" + url)
	return err
}

func (s Service) MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]entity.File, error) {
	var links []entity.File

	for _, f := range files {
		link, err := s.Upload(ctx, f, folder)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}
