package logic

import (
	"github.com/Alexchent/goscan/config"
	"github.com/mitchellh/go-homedir"
	"testing"
)

func Test_scanFile_Insert(t *testing.T) {
	dir, _ := homedir.Dir()
	fullFileName := dir + "/Downloads/No.6914梦心玥.zip"
	NewSaveLogic(config.Config{Sqlite: struct{ DSN string }{DSN: dir}}).Save(fullFileName)
}
