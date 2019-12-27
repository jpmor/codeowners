package codeowners

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func BuildEntriesFromFile(filePath string, includeComments bool) []*Entry {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	entries := []*Entry{}
	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}
		if len(line) < 1 {
			continue
		}

		parser := NewParser(strings.NewReader(string(line)))
		entry, err := parser.Parse()
		if err != nil {
			panic(err)
		}

		if (entry.suffix == PathSufix(None)) && !includeComments {
			continue
		}

		entries = append(entries, entry)

	}
	return entries
}
