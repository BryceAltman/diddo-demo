package utils

import (
	"fmt"
	"mime/multipart"
	"strings"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func ValidateImageFile(file *multipart.FileHeader) error {
	if file.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of %d bytes", MaxFileSize)
	}

	contentType := file.Header.Get("Content-Type")
	if !AllowedImageTypes[contentType] {
		return fmt.Errorf("unsupported file type: %s", contentType)
	}

	return nil
}

func GetFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}
	return strings.ToLower(parts[len(parts)-1])
}

func IsValidImageExtension(filename string) bool {
	ext := GetFileExtension(filename)
	validExtensions := []string{"jpg", "jpeg", "png", "gif", "webp"}
	
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}