package duplicate

import (
	myFile "github.com/Alexchent/goscan/file"
	"os"
	"sync"
)

const (
	consumer = 100
)

type File struct {
	FullFileName string
	MD5          string
}

type FIleList map[string][]string

func Do(dir string, B chan *File) {
	// 将文件信息写入channel A
	A := make(chan string, consumer)
	//wg.Add(1)

	var pWg sync.WaitGroup
	pWg.Add(1)
	go WalkFiles(dir, &pWg, A)

	// 从channel A 读取文件信息，并计算出MD5，以struct形式写入channel B
	go func() {
		// 当A关闭时退出
		for filename := range A {
			//fmt.Println(filename)
			B <- &File{
				FullFileName: filename,
				MD5:          myFile.GetFileMd5(filename),
			}
		}
		close(B)
	}()

	pWg.Wait()
	close(A)
	return
}

func WalkFiles(dir string, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()
	readDir, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var nextWg sync.WaitGroup
	for _, file := range readDir {
		fileName := file.Name()
		// 过滤影藏文件
		if fileName[0] == '.' {
			continue
		}
		// 按后缀过滤
		fullName := dir + "/" + fileName
		if file.IsDir() {
			nextWg.Add(1)
			go WalkFiles(fullName, &nextWg, ch)
		} else {
			ch <- fullName
		}
	}
	// 等待子goroutine结束。关闭channel，通知消费者退出
	nextWg.Wait()
	return
}
