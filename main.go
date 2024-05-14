package main

import (
	"fmt"
	"toolbox/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("failed to run toolbox:", err.Error())
	}
}
