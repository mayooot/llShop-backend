package check

import (
	"errors"
	"mime/multipart"
	"path/filepath"
)

// 允许的图片后缀名
var allowExts = []string{
	".jpg",
	".png",
	".gif",
	".jpeg",
}

var ErrorFileExtError = errors.New("不允许的图片后缀名")

// CheckPic 检查图片大小和格式
func CheckPic(fileHeader *multipart.FileHeader) (err error) {
	// todo 判断文件大小
	fileExt := filepath.Ext(fileHeader.Filename)
	allowFlag := false
	for _, ext := range allowExts {
		if ext == fileExt {
			allowFlag = true
			break
		}
	}
	if !allowFlag {
		err = ErrorFileExtError
		return
	}
	return
}
