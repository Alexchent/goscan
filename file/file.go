package myFile

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// CreateDateDir basePath是固定目录路径
func CreateDateDir(folderPath string) (dirPath string) {
	if folderPath == "" {
		return ""
	}
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步
		// 先创建文件夹
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		// 再修改权限
		err = os.Chmod(folderPath, 0777)
		if err != nil {
			log.Fatal(err)
			return ""
		}
	}
	return folderPath
}

// ReadString 逐行读取文件
func ReadString(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && io.EOF == err {
			break
		}
		data := strings.Trim(line, "\n")

		if data != "" {
			fmt.Println(data)
		}
	}
}

// ReadLine 逐行读取文件
func ReadLine(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
	}
}

// 生成文件的md5
func GetFileMd5(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 获取文件的md5
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%x\n", h.Sum(nil))
	return fmt.Sprintf("%x", h.Sum(nil))

	// 获取文件的sha1
	//h := sha1.New()
	//if _, err := io.Copy(h, file); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%x\n", h.Sum(nil))

	// 获取文件的sha256
	//h := sha256.New()
	//if _, err := io.Copy(h, file); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%x\n", h.Sum(nil))

	// 获取文件的sha512
	//h := sha512.New()
	//if _, err := io.Copy(h, file); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%x\n", h.Sum(nil))

	return ""
}
