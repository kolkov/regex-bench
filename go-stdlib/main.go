package main

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

type Pattern struct {
	Name    string
	Pattern string
}

var patterns = []Pattern{
	{"literal_alt", `error|warning|fatal|critical`},
	{"anchored", `^HTTP/[12]\.[01]`},
	{"inner_literal", `.*@example\.com`},
	{"suffix", `.*\.(txt|log|md)`},
	{"char_class", `[\w]+`},
	{"email", `[\w.+-]+@[\w.-]+\.[\w.-]+`},
}

func measure(data []byte, p Pattern) {
	start := time.Now()

	re := regexp.MustCompile(p.Pattern)
	matches := re.FindAll(data, -1)
	count := len(matches)

	elapsed := time.Since(start)
	ms := float64(elapsed) / float64(time.Millisecond)

	fmt.Printf("%-15s %10.2f ms  %6d matches\n", p.Name, ms, count)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: benchmark <input-file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Go stdlib (input: %.2f MB)\n", float64(len(data))/1024/1024)
	fmt.Println("─────────────────────────────────────────")

	for _, p := range patterns {
		measure(data, p)
	}
}
