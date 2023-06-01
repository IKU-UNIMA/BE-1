package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func CheckFileIsExcel(file *multipart.FileHeader) error {
	fileType := file.Header["Content-Type"][0]
	if fileType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		return nil
	}

	return FailedResponse(http.StatusBadRequest, map[string]string{"message": "unsupported file type for " + file.Filename})
}

func WriteFile(file *multipart.FileHeader, fileName string) error {
	if err := CheckFileIsExcel(file); err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return FailedResponse(http.StatusInternalServerError, nil)
	}
	defer src.Close()

	dst, err := os.Create(fileName)
	if err != nil {
		return FailedResponse(http.StatusInternalServerError, nil)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return FailedResponse(http.StatusInternalServerError, nil)
	}

	return nil
}

func GetNewFileName(name string) string {
	fileName := ""
	ext := ""
	for i := 0; i < len(name); i++ {
		if string(name[len(name)-i-1]) == "." {
			fileName = name[:len(name)-i-1]
			ext = name[len(name)-i-1:]
		}
	}

	return fmt.Sprintf("%s%d%s", fileName, time.Now().Unix(), ext)
}
