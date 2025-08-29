package logic

import (
	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/model"
)

type SaveLogic struct {
	Model *model.ScanFile
}

func NewSaveLogic(config config.Config) *SaveLogic {
	return &SaveLogic{Model: model.NewScanFile(config.Dir)}
}

func (s *SaveLogic) Save(fullFileName string) {
	s.Model.Insert(fullFileName)
}

func (s *SaveLogic) FIndByFullFileName(fullFileName string) *model.MyFile {
	return s.Model.FindByFullFileName(fullFileName)
}
