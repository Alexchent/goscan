package help

import (
	"fmt"
	"os"
)

func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	} else {
		fmt.Println(err.Error())
		return true
	}
}
