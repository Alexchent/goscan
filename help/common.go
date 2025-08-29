package help

import (
	"crypto/md5"
	"fmt"
)

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func FormatFileSize(fileSize int64) string {
	if fileSize < 1024 {
		return fmt.Sprintf("%d B", fileSize)
	} else if fileSize < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(fileSize)/1024)
	} else if fileSize < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(fileSize)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(fileSize)/(1024*1024*1024))
	}
}
