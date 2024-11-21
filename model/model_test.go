package model

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"testing"
)

func Test_scanFile_Insert(t *testing.T) {
	dir, _ := homedir.Dir()
	fullFileName := dir + "/Downloads/No.6914梦心玥.zip"
	NewScanFile("").Insert(fullFileName)
}

func Test_scanFile_Find(t *testing.T) {
	dir, _ := homedir.Dir()
	fullFileName := dir + "/Downloads/No.6914梦心玥.zip"
	t.Log(fullFileName)
	data := NewScanFile("").FindByFullFileName(fullFileName)
	fmt.Println("fileMd5: " + data.FileMD5)
	fmt.Println(fmt.Sprintf("fileSize: %fkb", float64(data.Size)/1024/1024))
}

func Test_scanFile_Find_Md5(t *testing.T) {
	fileMd5 := "83409ae895629b29af99d82090753004"
	data := NewScanFile("").FindMd5(fileMd5)
	if len(data) > 0 {
		for _, item := range data {
			fmt.Println("filename: " + item.FullFileName)
			fmt.Println(fmt.Sprintf("fileSize: %fkb", float64(item.Size)/1024/1024))
		}
	}
}
