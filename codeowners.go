package codeowners

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func BuildEntriesFromFile(filePath string) []*Entry {
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
		log.Println(entry)

		entries = append(entries, entry)

	}
	return entries
}