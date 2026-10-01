package main

import (
	"bufio"
	"fmt"
	"io"
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
		// fields, _ := shlex.Split(cmd)

		if len(fields) == 0 {
			os.Exit(0)
		}

		var out, errOut io.Writer = os.Stdout, os.Stderr
		var files []*os.File
		var args []string
		failed := false

		for i := 0; i < len(fields); i++ {
			op := fields[i]
			if (op == ">" || op == "1>" || op == "2>" || op == ">>" || op == "1>>" || op == "2>>") && i+1 < len(fields) {
				flag := os.O_TRUNC
				if op == ">>" || op == "1>>" || op == "2>>" {
					flag = os.O_APPEND
				}
				f, err := os.OpenFile(fields[i+1], os.O_CREATE|os.O_WRONLY|flag, 0644)
				if err != nil {
					fmt.Fprintln(errOut, err)
					failed = true
					break
				}
				files = append(files, f)
				if op == "2>" || op == "2>>" {
					errOut = f
				} else {
					out = f
				}
				i++
				continue
			}
			args = append(args, op)
		}

		if !failed && len(args) > 0 {
			if slices.Contains(builtins, args[0]) {
				handleCommand(args, out, errOut)
			} else if _, err := exec.LookPath(args[0]); err == nil {
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = out
				cmd.Stderr = errOut
				cmd.Run()
			} else {
				fmt.Fprintf(errOut, "%s: command not found\n", args[0])
			}
		}

		for _, f := range files {
			f.Close()
		}
	}
}

// alternative: use shlex library to parse
func parseFields(cmd string) []string {
	var fields []string
	var field strings.Builder
	var quote rune
	inField := false
	escaped := false

	for _, r := range cmd {
		if escaped {
			field.WriteRune(r)
			inField = true
			escaped = false
			continue
		}
		if quote != 0 {
			if quote == '"' && r == '\\' {
				escaped = true
				continue
			}
			if r == quote {
				quote = 0
			} else {
				field.WriteRune(r)
			}
			inField = true
			continue
		}
		if r == '\\' {
			escaped = true
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

	if escaped {
		field.WriteRune('\\')
	}

	if inField {
		fields = append(fields, field.String())
	}

	return fields
}

func handleType(cmd string, out, errOut io.Writer) {
	if slices.Contains(builtins, cmd) {
		fmt.Fprintf(out, "%s is a shell builtin\n", cmd)
		return
	} else if path, err := exec.LookPath(cmd); err == nil {
		fmt.Fprintf(out, "%s is %s\n", cmd, path)
		return
	}
	fmt.Fprintf(errOut, "%s: not found\n", cmd)
}

func handleCommand(fields []string, out, errOut io.Writer) {
	switch fields[0] {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Fprintln(out, strings.Join(fields[1:], " "))
	case "type":
		handleType(fields[1], out, errOut)
	case "pwd":
		wd, _ := os.Getwd()
		fmt.Fprintln(out, wd)
	case "cd":
		path := fields[1]
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
		if err := os.Chdir(path); err != nil {
			fmt.Fprintf(errOut, "cd: %s: No such file or directory\n", path)
		}
	}
}
