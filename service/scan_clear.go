package main

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	"strings"
)

// 清理后缀是js和torrent的文件
func main() {
	key := "have_save_file"
	val := mredis.SMembers(key)
	for _, v := range val {
		newv := strings.TrimRight(v, "\n")
		mredis.SRem(key, v)
		mredis.SAdd(key, newv)

		// 按 后缀清理
		if strings.HasSuffix(v, "js") || strings.HasSuffix(v, "torrent") {
			fmt.Println(v)
			mredis.SRem(key, v)
		}
	}
}
