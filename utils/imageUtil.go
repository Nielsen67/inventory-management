package utils

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

var validExtension = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func IsValidExtension(image *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(image.Filename))
	return validExtension[ext]
}
