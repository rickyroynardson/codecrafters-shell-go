package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		i, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			os.Exit(1)
		}

		i = strings.TrimSpace(i)
		inputs := strings.Split(i, " ")
		cmd := inputs[0]
		args := inputs[1:]
		switch cmd {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}
