package main

import (
	"fmt"
	"github.com/aeof/toolbox/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("failed to run toolbox:", err.Error())
	}
}
