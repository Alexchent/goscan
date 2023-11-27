package service

import (
	mconf "github.com/Alexchent/goscan/config"
	myFile "github.com/Alexchent/goscan/file"
	"log"
	"os"
	"strings"
)

// FileList map的key是文件的md5，value是文件的路径
var FileList = make(map[string][]string)

func Search(filePath string) {
	fileInfoList, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for i := range fileInfoList {
		fileName := fileInfoList[i].Name()
		if fileInfoList[i].IsDir() {
			WriteToFile(filePath + "/" + fileName)
		} else {
			if fileInfoList[i].Name() == ".DS_Store" {
				continue
			}
			// 判断是否是忽略的文件类型
			ignore := false
			for _, v := range mconf.Conf.FilterType {
				if strings.HasSuffix(fileName, v) {
					//log.Println("忽略文件：", fileName)
					ignore = true
				}
			}
			if ignore {
				continue
			}

			filename := filePath + "/" + fileName
			fileMd5 := myFile.GetFileMd5(filename)
			// 判断文件是否存在
			if _, ok := FileList[fileMd5]; ok {
				FileList[fileMd5] = append(FileList[fileMd5], filename)
			} else {
				FileList[fileMd5] = []string{filename}
			}
			// 保存到redis成功，说明是新的文件
			//if mredis.SAdd(CacheKeyMd5, fileMd5) == 1 {
			//	fmt.Println("发现新的文件：", filename)
			//	myFile.AppendContent("history_file.txt", fileMd5+"/t"+filename)
			//} else {
			//	myFile.AppendContent("same_file.txt", fileMd5+"/t"+filename)
			//}
		}
	}
}
