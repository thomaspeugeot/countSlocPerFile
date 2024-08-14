package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

type FuncSLOC struct {
	Name string
	SLOC int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <go-source-file>")
		return
	}

	filename := os.Args[1]
	fs := token.NewFileSet()

	// Parse the Go source file
	node, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Read file lines into a slice
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Traverse the AST and count SLOC per function
	var funcs []FuncSLOC
	ast.Inspect(node, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			start := fs.Position(fn.Pos()).Line - 1
			end := fs.Position(fn.End()).Line - 1
			sloc := countSLOC(lines[start:end])

			funcs = append(funcs, FuncSLOC{Name: fn.Name.Name, SLOC: sloc})
		}
		return true
	})

	// Sort functions by ascending SLOC
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].SLOC < funcs[j].SLOC
	})

	// Print sorted results
	for _, fn := range funcs {
		fmt.Printf("Function %s has %d SLOC\n", fn.Name, fn.SLOC)
	}
}

func countSLOC(lines []string) int {
	count := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "/*") && !strings.HasPrefix(trimmed, "*") && !strings.HasPrefix(trimmed, "*/") {
			count++
		}
	}
	return count
}
