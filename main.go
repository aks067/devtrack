package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: devtrack [daemon|status]")
		return
	}
	if os.Args[1] == "daemon" {
		fmt.Println("mode daemon")
	} else if os.Args[1] == "status" {
		fmt.Println("mode status")
	}
}
