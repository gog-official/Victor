package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gog-official/victor/internal/lexer"
	"github.com/gog-official/victor/internal/parser"
	"github.com/gog-official/victor/internal/vector"
	"github.com/gog-official/victor/internal/workspace"
)

var (
	heading string = "Victor does vector"
	name    string = "Victor"
	store          = workspace.NewWorkspace()
)

func printPrompt() {
	fmt.Print(name, "> ")
}

func printHeader() {
	fmt.Println(heading)
}

func handleHelp() {
	fmt.Printf("Welcome to %v! These are available commands: \n", name)
	fmt.Println()
	fmt.Println("Expressions:")
	fmt.Println("  1/3 + 1/6          - Scalar arithmetic")
	fmt.Println("  a + b              - Vector addition")
	fmt.Println("  a - b              - Vector subtraction")
	fmt.Println("  a * 2  or  2 * a   - Scaling")
	fmt.Println("  a / 2              - Dividing a vector by a number")
	fmt.Println("  a * b              - Dot product")
	fmt.Println("  -a                 - Negation")
	fmt.Println()
	fmt.Println("Vectors (maximum 5):")
	fmt.Println("  define a = (3, 4) - Defines or replaces vector a")
	fmt.Println("  define c = a + b  - Defines a vector from an expression")
	fmt.Println("  show a            - Shows vector a")
	fmt.Println("  list              - Lists all vectors")
	fmt.Println("  delete a          - Deletes vector a")
	fmt.Println()
	fmt.Println("Vector functions:")
	fmt.Println("  mag a             - Magnitude of a")
	fmt.Println("  norm a            - Normalizes a")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help              - Shows available operations")
	fmt.Println("  clear             - Clears the terminal")
	fmt.Println("  exit              - Closes Victor. please dont :( ")
}

func handleClear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func evaluate(txt string) (parser.Value, error) {
	lex := lexer.NewLexer(txt)
	p := parser.NewParser(lex, store)
	return p.ParseValue()
}

func formatVector(v vector.Vector2D) string {
	return fmt.Sprintf("(%s, %s)", v.X.String(), v.Y.String())
}

func printValue(value parser.Value) {
	if value.IsVector {
		fmt.Println(formatVector(value.Vector))
		return
	}

	fmt.Println(value.Scalar.String())
}

func handleDefine(txt string) error {
	body := strings.TrimSpace(strings.TrimPrefix(txt, "define"))
	parts := strings.SplitN(body, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("usage: define <name> = (<x>, <y>)")
	}

	vecName := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if value == "" {
		return fmt.Errorf("usage: define <name> = (<x>, <y>)")
	}

	result, err := evaluate(value)
	if err != nil {
		return err
	}

	if !result.IsVector {
		return fmt.Errorf("expected a vector")
	}

	if err := store.Define(vecName, result.Vector); err != nil {
		return err
	}

	fmt.Printf("%s = %s\n", workspace.FixName(vecName), formatVector(result.Vector))
	return nil
}

func handleShow(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: show <name>")
	}

	vec, exists := store.Get(args[1])
	if !exists {
		return fmt.Errorf("vector %q is not defined", args[1])
	}

	fmt.Printf("%s = %s\n", workspace.FixName(args[1]), formatVector(vec))
	return nil
}

func handleList() {
	names := store.List()
	if len(names) == 0 {
		fmt.Println("no vectors defined")
		return
	}

	for _, vecName := range names {
		vec, _ := store.Get(vecName)
		fmt.Printf("%s = %s\n", vecName, formatVector(vec))
	}
}

func handleDelete(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: delete <name>")
	}

	if !store.Delete(args[1]) {
		return fmt.Errorf("vector %q is not defined", args[1])
	}

	fmt.Printf("deleted %s\n", workspace.FixName(args[1]))
	return nil
}

func handleMagnitude(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: mag <vector>")
	}

	result, err := evaluate(strings.Join(args[1:], " "))
	if err != nil {
		return err
	}

	if !result.IsVector {
		return fmt.Errorf("expected a vector")
	}

	fmt.Println(result.Vector.Magnitude())
	return nil
}

func handleNormalize(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: norm <vector>")
	}

	result, err := evaluate(strings.Join(args[1:], " "))
	if err != nil {
		return err
	}

	if !result.IsVector {
		return fmt.Errorf("expected a vector")
	}

	normalized, err := result.Vector.Normalize()
	if err != nil {
		return err
	}

	fmt.Println(formatVector(normalized))
	return nil
}

func handleCmd(txt string) {
	if txt == "" {
		return
	}

	args := strings.Fields(txt)
	command := args[0]

	switch command {
	case "define":
		if err := handleDefine(txt); err != nil {
			fmt.Println("error:", err)
		}
	case "show":
		if err := handleShow(args); err != nil {
			fmt.Println("error:", err)
		}
	case "list":
		handleList()
	case "delete":
		if err := handleDelete(args); err != nil {
			fmt.Println("error:", err)
		}
	case "mag", "magnitude":
		if err := handleMagnitude(args); err != nil {
			fmt.Println("error:", err)
		}
	case "norm", "normalize":
		if err := handleNormalize(args); err != nil {
			fmt.Println("error:", err)
		}
	default:
		result, err := evaluate(txt)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		printValue(result)
	}
}

func cleanInput(txt string) string {
	output := strings.TrimSpace(txt)
	return strings.ToLower(output)
}

func main() {
	commands := map[string]interface{}{
		"help":  handleHelp,
		"clear": handleClear,
	}

	scanner := bufio.NewScanner(os.Stdin)
	printHeader()
	printPrompt()

	for scanner.Scan() {
		txt := cleanInput(scanner.Text())
		if command, exists := commands[txt]; exists {
			command.(func())()
		} else if strings.EqualFold("exit", txt) {
			return
		} else {
			handleCmd(txt)
		}

		printPrompt()
	}
	fmt.Println()
}
