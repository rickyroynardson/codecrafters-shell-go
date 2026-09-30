package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	builtins := []string{"echo", "exit", "type"}

	for {
		fmt.Print("$ ")

		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			os.Exit(1)
		}

		cmd = strings.TrimSpace(cmd)
		tokens := strings.Split(cmd, " ")

		switch tokens[0] {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(tokens[1:], " "))
		case "type":
			// panic: need to handle index out of range if args is empty
			if slices.Contains(builtins, tokens[1]) {
				fmt.Printf("%s is a shell builtin\n", tokens[1])
			} else {
				path, err := exec.LookPath(tokens[1])
				if err != nil || path == "" {
					fmt.Printf("%s: not found\n", tokens[1])
				} else {
					fmt.Printf("%s is %s\n", tokens[1], path)
				}
			}
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}
