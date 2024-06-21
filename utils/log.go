package utils

import "fmt"

var Verbose bool

func LogVerbose(msg ...any) {
	if Verbose {
		fmt.Println(msg...)
	}
}
