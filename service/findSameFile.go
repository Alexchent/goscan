package service

import (
	"fmt"
	mconf "github.com/Alexchent/goscan/config"
	myFile "github.com/Alexchent/goscan/file"
	"log"
	"os"
	"os/exec"
	"path"
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
		if fileName == ".git" {
			continue
		}
		if fileInfoList[i].IsDir() {
			Search(filePath + "/" + fileName)
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
			fmt.Println(filename)
			// 判断文件是否存在
			if _, ok := FileList[fileMd5]; ok {
				FileList[fileMd5] = append(FileList[fileMd5], filename)
			} else {
				FileList[fileMd5] = []string{filename}
			}
		}
	}
}

func GetSame() {
	//fmt.Println(FileList)
	same, _ := os.OpenFile("same_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	//history, _ := os.OpenFile("history_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	for k, v := range FileList {
		if len(v) > 1 {
			//fmt.Println("发现相同文件：", k)
			//myFile.AppendContent("same_file.txt", k)
			same.WriteString(k + "\n")
			for _, v := range v {
				same.WriteString("\t" + v + "\n")
				//myFile.AppendContent("same_file.txt", "\t"+v)
			}
		}
		//for _, v := range v {
		//	history.WriteString(k + "\t" + v + "\n")
		//}
	}
}

func RemoveSameFile() {
	for _, v := range FileList {
		if len(v) > 1 {
			// 保留第一个文件，删除其他文件
			for _, v := range v[1:] {
				fmt.Println("remove:\t" + v)
				cmd := exec.Command("mv", v, "/Users/chentao/same/"+path.Base(v))
				err := cmd.Run()
				if err != nil {
					log.Printf("rename %s failed: %v", v, err)
				}
			}
		}
	}
}
