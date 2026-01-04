//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const (
	targetSize = 6 * 1024 * 1024 // 6 MB
)

var (
	words = []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
		"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
		"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
		"function", "return", "if", "else", "while", "for", "var", "const", "let",
		"import", "export", "class", "interface", "type", "struct", "package",
		"server", "client", "request", "response", "data", "file", "config",
	}

	// Samples that will match our patterns
	logLevels    = []string{"error", "warning", "fatal", "critical"}
	httpVersions = []string{"HTTP/1.0", "HTTP/1.1", "HTTP/2.0"}
	emails       = []string{"user@example.com", "admin@test.org", "info@company.net", "support@example.com"}
	filenames    = []string{"readme.txt", "config.log", "notes.md", "data.txt", "server.log", "docs.md"}
	uris         = []string{
		"http://example.com/path/to/resource",
		"https://api.github.com/repos/user/repo?page=1",
		"ftp://files.server.net/downloads/file.zip",
		"https://www.google.com/search?q=regex#results",
	}
	ips = []string{
		"192.168.1.1", "10.0.0.255", "172.16.0.1", "255.255.255.0",
		"8.8.8.8", "127.0.0.1", "203.0.113.42", "198.51.100.7",
	}

	// Fruits for multi-literal Aho-Corasick test (>8 patterns)
	fruits = []string{
		"apple", "banana", "cherry", "date", "elderberry", "fig",
		"grape", "honeydew", "kiwi", "lemon", "mango", "orange",
	}

	// Version numbers for digit-start DigitPrefilter test
	versions = []string{
		"1.0.0", "2.1.3", "3.14.159", "10.20.30", "0.9.1", "4.5.6",
		"1.2.3", "7.8.9", "12.0.1", "2.0.0", "5.4.3", "99.99.99",
	}
)

func randomWord() string {
	return words[rand.Intn(len(words))]
}

func randomWords(n int) string {
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = randomWord()
	}
	return strings.Join(parts, " ")
}

func generateContent() string {
	var builder strings.Builder
	builder.Grow(targetSize + 1024)

	lineNum := 0
	for builder.Len() < targetSize {
		lineNum++
		var line string

		// Every 100 lines, add a special line with guaranteed matches
		switch {
		case lineNum%500 == 1:
			// HTTP request line (anchored pattern)
			line = fmt.Sprintf("%s 200 OK %s", httpVersions[rand.Intn(len(httpVersions))], randomWords(5))

		case lineNum%100 == 2:
			// Log line with level
			level := logLevels[rand.Intn(len(logLevels))]
			line = fmt.Sprintf("[%s] %s %s", level, randomWords(8), filenames[rand.Intn(len(filenames))])

		case lineNum%150 == 3:
			// Email line
			email := emails[rand.Intn(len(emails))]
			line = fmt.Sprintf("Contact: %s for %s", email, randomWords(6))

		case lineNum%80 == 4:
			// Filename line
			filename := filenames[rand.Intn(len(filenames))]
			line = fmt.Sprintf("File: %s - %s", filename, randomWords(7))

		case lineNum%120 == 5:
			// URI line
			uri := uris[rand.Intn(len(uris))]
			line = fmt.Sprintf("Link: %s %s", uri, randomWords(5))

		case lineNum%90 == 6:
			// IP address line
			ip := ips[rand.Intn(len(ips))]
			line = fmt.Sprintf("Server %s responded with %s", ip, randomWords(4))

		case lineNum%70 == 7:
			// Fruit line (for multi-literal Aho-Corasick test)
			fruit := fruits[rand.Intn(len(fruits))]
			line = fmt.Sprintf("Fresh %s available at the market %s", fruit, randomWords(5))

		case lineNum%60 == 8:
			// Version line (for digit-start DigitPrefilter test)
			version := versions[rand.Intn(len(versions))]
			line = fmt.Sprintf("Updated to version %s with %s", version, randomWords(6))

		default:
			// Regular line with random content
			wordCount := 8 + rand.Intn(12)
			line = randomWords(wordCount)
		}

		builder.WriteString(line)
		builder.WriteByte('\n')
	}

	return builder.String()
}

func main() {
	rand.Seed(42) // Fixed seed for reproducibility

	data := generateContent()

	if err := os.MkdirAll("input", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating dir: %v\n", err)
		os.Exit(1)
	}

	err := os.WriteFile("input/data.txt", []byte(data), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated input/data.txt (%.2f MB)\n", float64(len(data))/1024/1024)
}
