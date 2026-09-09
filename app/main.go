package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("$ ")
	i, _ := reader.ReadString('\n')
	fmt.Printf("%s: command not found\n", i[:len(i)-1])
}
