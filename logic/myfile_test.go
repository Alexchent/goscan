package logic

import (
	"github.com/Alexchent/goscan/config"
	"github.com/mitchellh/go-homedir"
	"testing"
	"time"
)

func Test_scanFile_Insert(t *testing.T) {
	dir, _ := homedir.Dir()
	fullFileName := dir + "/Downloads/No.6914梦心玥.zip"
	NewSaveLogic(config.Config{Dir: dir}).Save(fullFileName)
}

func Test_scanFile_Find(t *testing.T) {
	dir, _ := homedir.Dir()
	fullFileName := dir + "/Downloads/No.6914梦心玥.zip"
	data := NewSaveLogic(config.Config{Dir: dir}).FIndByFullFileName(fullFileName)
	t.Log("fullFileName: " + data.FullFileName)
	t.Log("fileMd5: " + data.FileMD5)
	t.Logf("fileSize: %d", data.Size)
	t.Logf("created:" + data.Created.Format(time.DateTime))
}
