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
	builtins := []string{"echo", "exit", "type", "pwd"}

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
			if len(tokens) < 2 {
				fmt.Println("type: missing argument")
				continue
			}
			target := tokens[1]
			if slices.Contains(builtins, target) {
				fmt.Printf("%s is a shell builtin\n", target)
			} else if path, err := exec.LookPath(target); err == nil {
				fmt.Printf("%s is %s\n", target, path)
			} else {
				fmt.Printf("%s: not found\n", target)
			}
		case "pwd":
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("error getting current directory: %v\n", err)
			}
			fmt.Printf("%s\n", wd)
		default:
			_, err := exec.LookPath(tokens[0])
			if err == nil {
				cmd := exec.Command(tokens[0], tokens[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Printf("error executing command: %v\n", err)
				}
			} else {
				fmt.Printf("%s: command not found\n", cmd)
			}
		}
	}
}
