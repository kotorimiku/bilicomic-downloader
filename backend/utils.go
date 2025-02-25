package bilicomicdownloader

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CreateZipFromDirectory(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if fi.IsDir() {
			return nil
		}

		// 创建一个文件在压缩包中
		relPath, err := filepath.Rel(sourceDir, file)
		if err != nil {
			return err
		}

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileReader, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		_, err = io.Copy(writer, fileReader)
		return err
	})

	return err
}

var illegalChars = regexp.MustCompile(`[<>:"/\\|?*]+`)

func sanitizeFilename(filename string) string {
	filename = illegalChars.ReplaceAllString(filename, "")

	filename = strings.TrimRight(filename, ".")

	if filename == "" {
		filename = "unnamed"
	}

	return filename
}

func isImage(data []byte) bool {
	// PNG文件的前缀字节
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return true
	}

	// JPEG文件的前缀字节
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return true
	}

	// GIF文件的前缀字节
	if bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46, 0x38}) {
		return true
	}

	// WebP文件的前缀字节
	if bytes.HasPrefix(data, []byte{'R', 'I', 'F', 'F'}) && bytes.HasPrefix(data[8:], []byte{'W', 'E', 'B', 'P'}) {
		return true
	}

	// AVIF文件的前缀字节
	if bytes.HasPrefix(data[4:], []byte("ftypavif")) {
		return true
	}

	// 其他常见格式可以继续扩展
	return false
}
