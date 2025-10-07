package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go-rest-setup/src/database/models"
	config "go-rest-setup/src/lib/configs"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileService struct {
	config *config.Config
	db     *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{
		config: config.EnvModule(),
		db:     db,
	}
}

func (s *FileService) validateFile(file *multipart.FileHeader) error {
	maxSizeMB := 5 // MB
	if file.Size > int64(maxSizeMB*1024*1024) {
		return fmt.Errorf("file size exceeds %dMB", maxSizeMB)
	}

	allowed := []string{
		"image/jpeg", "image/png", "image/webp", "image/gif",
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}
	for _, t := range allowed {
		if file.Header.Get("Content-Type") == t {
			return nil
		}
	}
	return errors.New("unsupported file type")
}

func (s *FileService) generateFileName(original string) string {
	ext := filepath.Ext(original)
	base := strings.TrimSuffix(filepath.Base(original), ext)
	base = strings.ReplaceAll(base, " ", "_")
	return fmt.Sprintf("%s_%s%s", base, uuid.New().String()[:6], ext)
}

func (s *FileService) resizeImage(file multipart.File, width, height int) (*bytes.Buffer, error) {
	img, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}
	resized := imaging.Resize(img, width, height, imaging.Lanczos)
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, resized, imaging.JPEG)
	return buf, err
}

func (s *FileService) Upload(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*models.File, error) {
	if err := s.validateFile(fileHeader); err != nil {
		return nil, err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	fileName := s.generateFileName(fileHeader.Filename)

	var key string
	if folder == "" || folder == "/" {
		key = fileName
	} else {
		key = fmt.Sprintf("%s/%s", strings.Trim(folder, "/"), fileName)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		return nil, err
	}

	_, err = config.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &[]string{s.config.S3.Bucket}[0],
		Key:         &key,
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: &[]string{fileHeader.Header.Get("Content-Type")}[0],
		// ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/%s",
		s.config.S3.URL,
		s.config.S3.Bucket,
		key,
	)

	file := &models.File{
		Nama:         fileName,
		OriginalName: fileHeader.Filename,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Size:         fileHeader.Size,
		URL:          url,
	}

	if err := s.db.Create(&file).Error; err != nil {
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	return file, nil
}
