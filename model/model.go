package model

import (
	"github.com/Alexchent/goscan/help"
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

type ScanFile interface {
	Insert(fullFileName string)
	FindByFullFileName(fullFileName string) *MyFile
	FindMd5(fileMd5 string) []MyFile
}

type scanFile struct {
	db *gorm.DB
}

// 接口完整性校验
var _ ScanFile = (*scanFile)(nil)

func NewScanFile(dsn string) *scanFile {
	if len(dsn) > 0 && dsn[len(dsn)-1] != '/' {
		dsn += "/"
	}
	db, err := gorm.Open(sqlite.Open(dsn+"scan.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 迁移 schema
	err = db.AutoMigrate(&MyFile{})
	if err != nil {
		panic(err)
		return nil
	}
	return &scanFile{db: db}
}

func (m *scanFile) Insert(fullFileName string) {
	var fileMD5 string
	var fileSize int64
	result := m.FindByFullFileName(fullFileName)
	stat, err := os.Stat(fullFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("文件 %s 不存在\n", fullFileName)
		}
		fileSize = 0
	} else {
		fileSize = stat.Size()
		fileMD5 = help.GetFileMd5(fullFileName)
	}
	if result.FullFileName == "" {
		m.db.Create(&MyFile{
			FullFileName: fullFileName,
			Size:         fileSize,
			FileMD5:      fileMD5,
			Created:      time.Now(),
			Status:       1,
		})
	} else if result.Size == 0 {
		result.Size = fileSize
		m.db.Save(result)
	}
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
