package model

import (
	myFile "github.com/Alexchent/goscan/file"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type MyFile struct {
	FullFileName string
	Size         int64 // 文件大小
	FileMD5      string
	Created      time.Time
	Status       int
}

var db *gorm.DB

func init() {
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	//if err != nil {
	//	panic("failed to connect database")
	//}
	//
	//// 迁移 schema
	//db.AutoMigrate(&MyFile{})
}

type ScanFile interface {
	Insert(fullFileName string, fileSize int64)
}

type scanFile struct {
	db *gorm.DB
}

func NewScanFile() *scanFile {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&MyFile{})
	return &scanFile{db: db}
}

func (m *scanFile) Insert(fullFileName string) {
	var fileMD5 string
	var fileSize int64
	stat, err := os.Stat(fullFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("文件 %s 不存在\n", fullFileName)
		}
		fileSize = 0
	} else {
		fileSize = stat.Size()
		fileMD5 = myFile.GetFileMd5(fullFileName)
	}
	m.db.Create(&MyFile{
		FullFileName: fullFileName,
		Size:         fileSize,
		FileMD5:      fileMD5,
		Created:      time.Now(),
		Status:       1,
	})
}

func (m *scanFile) FindByFullFileName(fullFileName string) *MyFile {
	var data MyFile
	m.db.First(&data, "full_file_name = ?", fullFileName)
	return &data
}

func (m *scanFile) FindMd5(fileMd5 string) []MyFile {
	var data []MyFile
	m.db.Find(&data, "file_md5 = ?", fileMd5)
	//m.db.Where("file_md5 = ?", fileMd5).Find(&data)
	return data
}
