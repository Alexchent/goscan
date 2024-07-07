package duplicate

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	consumer = 100
)

type File struct {
	FullFileName string
	MD5          string
}

type FIleList map[string][]string

func Do(dir string) {
	start := time.Now()
	defer fmt.Println("扫描完毕，共耗时：", time.Since(start))

	_, err := os.Stat(dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// stage 1
	A := make(chan string, consumer)
	wg := sync.WaitGroup{}
	for i := 0; i < consumer; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filename := range A {
				fmt.Println(filename)
			}
		}()
	}

	// 从channel A 读取文件信息，并计算出MD5，以struct形式写入channel B

	// 将文件信息写入channel A
	wg.Add(1)
	Stage1(dir, &wg, A)
	close(A)

	wg.Wait()
	fmt.Println("finish")
	// 从channel B 读取，写入到 map中
}

func Stage1(dir string, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()
	readDir, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	nextWg := sync.WaitGroup{}
	var fileName string
	for _, file := range readDir {
		fileName = file.Name()
		// 过滤影藏文件
		if fileName[0:1] == "." {
			continue
		}
		// 按后缀过滤
		fullName := dir + "/" + fileName
		if file.IsDir() {
			nextWg.Add(1)
			Stage1(fullName, &nextWg, ch)
		} else {
			ch <- fullName
		}
	}
	nextWg.Wait()
}
