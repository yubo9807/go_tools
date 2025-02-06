package utils

import "os"

func GetArgs(i int) string {
	for i < len(os.Args) {
		return os.Args[i]
	}
	return ""
}
