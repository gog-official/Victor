package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	heading string = "Victor does vector"
	name    string = "Victor"
)

func printPrompt() {
	fmt.Print(name, "> ")
}

func printHeader() {
	fmt.Println(heading)
}

func printUnknownCmd(txt string) {
	fmt.Println(txt, ": command not found")
}

func handleHelp() {
	fmt.Printf("Welcone to %v! These are available commandss: \n", name)
	fmt.Println(".help    - Shows available operations")
	fmt.Println(".clear   - Clears the terminal")
	fmt.Println(".exit    - Closes Victor. please dont :( ")
}

func handleClear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func handleUnknownCmd(txt string) {
	printUnknownCmd(txt)
}

func handleCmd(txt string) {
	handleUnknownCmd(txt)
}

func cleanInput(txt string) string {
	output := strings.TrimSpace(txt)
	return strings.ToLower(output)
}

func main() {
	commands := map[string]interface{}{
		".help":  handleHelp,
		".clear": handleClear,
	}

	scanner := bufio.NewScanner(os.Stdin)
	printHeader()
	printPrompt()

	for scanner.Scan() {
		txt := cleanInput(scanner.Text())
		if command, exists := commands[txt]; exists {
			command.(func())()
		} else if strings.EqualFold(".exit", txt) {
			return
		} else {
			handleCmd(txt)
		}

		printPrompt()
	}
	fmt.Println()
}
