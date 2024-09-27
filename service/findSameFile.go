package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	myFile "github.com/Alexchent/goscan/file"
	"log"
	"os"
	"os/exec"
	"path"
)

// FileList map的key是文件的md5，value是文件的路径
type FileList map[string]list
type list []string

// var listData = make(map[string][]string)
var listData = make(FileList)

func Search(filePath string) {
	fileInfoList, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for i := range fileInfoList {
		fileName := fileInfoList[i].Name()
		// 如果是影藏文件，直接跳过
		if fileName[0] == '.' {
			continue
		}
		if fileInfoList[i].IsDir() {
			Search(filePath + "/" + fileName)
		} else {
			filename := filePath + "/" + fileName
			fileMd5 := myFile.GetFileMd5(filename)
			fmt.Println(filename)
			// 判断文件是否存在
			if _, ok := listData[fileMd5]; ok {
				listData[fileMd5] = append(listData[fileMd5], filename)
			} else {
				listData[fileMd5] = []string{filename}
			}
		}
	}
}

func LogSameFile() FileList {
	//marshal, err := json.Marshal(listData)
	more := make(FileList)
	for k, v := range listData {
		if len(v) > 1 {
			more[k] = v
		}
	}
	if more == nil || len(more) == 0 {
		fmt.Println("没有重复文件")
		return nil
	}
	marshal, err := json.MarshalIndent(more, "", "     ")
	if err != nil {
		return nil
	}
	same, _ := os.OpenFile("same_file.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	_, _ = same.Write(marshal)
	return more
}

func RemoveSameFile(file FileList) {
	if len(file) == 0 {
		fmt.Println("没有重复文件")
		return
	}
	for _, v := range file {
		if len(v) > 1 {
			// 保留第一个文件，删除其他文件
			for _, v := range v[1:] {
				//fmt.Println("remove:\t" + v)
				encoded := MD5(v)
				cmd := exec.Command("mv", v, "/Users/chentao/same/"+encoded+"-"+path.Base(v))
				//err := cmd.Run()
				output, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("combined out:\n%s\n", string(output))
					//fmt.Printf("cmd.Run() failed with %s\n", err)
				}
			}
		}
	}
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
