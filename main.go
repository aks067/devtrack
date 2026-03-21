package main

import (
	"fmt"
	"github.com/aks067/devtrack/cmd"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: devtrack [daemon|status]")
		return
	}
	switch os.Args[1] {
	case "daemon":
		daemon.RunDaemon()
	case "status":
		//launch status func
	default:
		fmt.Println("Invalid Argument")
	}
}
