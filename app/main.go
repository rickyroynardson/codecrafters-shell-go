package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtins = []string{"echo", "exit", "type", "pwd", "cd"}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			os.Exit(1)
		}

		cmd = strings.TrimSpace(cmd)
		fields := parseFields(cmd)

		if len(fields) == 0 {
			os.Exit(0)
		}

		if slices.Contains(builtins, fields[0]) {
			handleCommand(fields)
		} else if _, err := exec.LookPath(fields[0]); err == nil {
			cmd := exec.Command(fields[0], fields[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			fmt.Printf("%s: command not found\n", fields[0])
		}
	}
}

// alternative: use shlex library to parse
func parseFields(cmd string) []string {
	var fields []string
	var field strings.Builder
	var quote rune
	inField := false

	for _, r := range cmd {
		if quote != 0 {
			if r == quote {
				quote = 0
			} else {
				field.WriteRune(r)
			}
			inField = true
			continue
		}
		if r == '\'' || r == '"' {
			quote = r
			inField = true
			continue
		}
		if r == ' ' || r == '\t' {
			if inField {
				fields = append(fields, field.String())
				field.Reset()
				inField = false
			}
			continue
		}
		field.WriteRune(r)
		inField = true
	}

	if inField {
		fields = append(fields, field.String())
	}

	return fields
}

func handleType(cmd string) {
	if slices.Contains(builtins, cmd) {
		fmt.Printf("%s is a shell builtin\n", cmd)
		return
	} else if path, err := exec.LookPath(cmd); err == nil {
		fmt.Printf("%s is %s\n", cmd, path)
		return
	}
	fmt.Printf("%s: not found\n", cmd)
}

func handleCommand(fields []string) {
	switch fields[0] {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(fields[1:], " "))
	case "type":
		handleType(fields[1])
	case "pwd":
		wd, _ := os.Getwd()
		fmt.Println(wd)
	case "cd":
		path := fields[1]
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
		if err := os.Chdir(path); err != nil {
			fmt.Printf("cd: %s: No such file or directory\n", path)
		}
	}
}
