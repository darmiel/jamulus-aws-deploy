package common

import (
	"os"
	"strings"
)

func SplitPath(s string) (dir, file string) {
	index := strings.LastIndex(s, string(os.PathSeparator))
	if index < 0 {
		return ".", ""
	}
	return s[:index], s[index+1:]
}
