package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
		builtins := map[string]struct{}{
			"echo": {},
			"exit": {},
			"type": {},
		}
		switch cmd {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			if _, ok := builtins[args[0]]; ok {
				fmt.Printf("%s is a shell builtin\n", args[0])
			} else {
				fmt.Printf("%s: not found\n", args[0])
			}
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}
