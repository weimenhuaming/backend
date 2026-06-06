package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxUploadSize   = 5 << 20 // 5MB
	StaticURLPrefix = "/static"
)

var staticRoot = "static"

// InitStaticRoot 设置静态文件根目录（需在服务启动时调用）
func InitStaticRoot(dir string) {
	if dir != "" {
		staticRoot = dir
	}
}

// StaticRoot 返回当前静态文件根目录
func StaticRoot() string {
	return staticRoot
}

var allowedImageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// SaveUploadedImage 保存 multipart 图片到 static/{subdir}/，返回可访问路径 /static/{subdir}/{filename}
func SaveUploadedImage(r *http.Request, fieldName, subdir string) (string, error) {
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		return "", fmt.Errorf("文件过大或格式无效")
	}

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("请选择要上传的图片")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExts[ext] {
		return "", fmt.Errorf("仅支持 JPG、PNG、GIF、WEBP 格式")
	}

	if err := validateImageContentType(header); err != nil {
		return "", err
	}

	dir := filepath.Join(staticRoot, subdir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("创建上传目录失败")
	}

	filename := uuid.New().String() + ext
	destPath := filepath.Join(dir, filename)

	dest, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("保存文件失败")
	}
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("保存文件失败")
	}

	return fmt.Sprintf("%s/%s/%s", StaticURLPrefix, subdir, filename), nil
}

func validateImageContentType(header *multipart.FileHeader) error {
	contentType := strings.ToLower(header.Header.Get("Content-Type"))
	if contentType == "" {
		return nil
	}

	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return nil
	default:
		return fmt.Errorf("仅支持 JPG、PNG、GIF、WEBP 格式")
	}
}
