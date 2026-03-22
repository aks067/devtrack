package main

import (
	"devtrack/cmd"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: devtrack [daemon|status]")
		return
	}
	switch os.Args[1] {
	case "daemon":
		cmds.RunDaemon()
	case "status":
		cmds.RunStatus()
	default:
		fmt.Println("Invalid Argument")
	}
}
