package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	// Read the README.md file
	content, err := ioutil.ReadFile("../README.md")
	if err != nil {
		log.Fatalf("Error reading README.md: %v", err)
	}

	// Regex to find migration entries
	re := regexp.MustCompile(`\*\s*\[.*\]\(.*\)\s*\(.*\)\s*from\s*(.*)\s*to\s*(.*)`)

	migrations := make(map[string]map[string]int)

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) == 3 {
			from := strings.TrimSpace(matches[1])
			to := strings.TrimSpace(matches[2])
			if migrations[from] == nil {
				migrations[from] = make(map[string]int)
			}
			migrations[from][to]++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning README.md: %v", err)
	}

	// Generate the .dot file content
	var dotFileContent strings.Builder
	dotFileContent.WriteString("digraph G {\n")
	dotFileContent.WriteString("  rankdir=LR;\n")
	dotFileContent.WriteString("  node [shape = record];\n")

	for from, toMap := range migrations {
		for to, count := range toMap {
			dotFileContent.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [penwidth=%d];\n", from, to, count))
		}
	}

	dotFileContent.WriteString("}\n")

	// Write the .dot file
	dotFilePath := "diagram.dot"
	err = ioutil.WriteFile(dotFilePath, []byte(dotFileContent.String()), 0644)
	if err != nil {
		log.Fatalf("Error writing dot file: %v", err)
	}

	// Generate the PNG image using dot
	outputImagePath := "../list-of-tech-migrations.png"
	cmd := exec.Command("dot", "-Tpng", dotFilePath, "-o", outputImagePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running dot command: %v", err)
	}

	fmt.Println("Diagram generated successfully!")
}
