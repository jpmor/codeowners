package codeowners

import (
	"bufio"
	"fmt"
	"path/filepath"

	"io"
	"log"
	"os"
	"strings"

	"github.com/alecharmon/trie"
)

type node struct {
	entries []*Entry
}

// CodeOwners search index for a CODEOWNER file
type CodeOwners struct {
	*trie.PathTrie
}

// BuildEntriesFromFile from an file path, absolute or relative, builds the entries for the CODEOWNERS file
func BuildEntriesFromFile(filePath string, includeComments bool) ([]*Entry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
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
	return entries, nil
}

func newNode() *node {
	return &node{
		entries: []*Entry{},
	}
}

// BuildFromFile from an file path, absolute or relative, builds the index for the CODEOWNERS file
func BuildFromFile(filePath string) (*CodeOwners, error) {
	t := &CodeOwners{
		trie.NewPathTrie(),
	}
	var n *node
	var ok bool

	entries, err := BuildEntriesFromFile(filePath, false)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		value := t.Get(entry.path)
		if value == nil {
			n = newNode()
		} else {
			n, ok = value.(*node)
			if !ok {
				out, _ := fmt.Printf("%v, %v", ok, n)
				panic(out)
			}
		}

		n.addEntry(entry)
		path := entry.path
		if []rune(path)[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		t.Put(path, n)
	}

	return t, nil
}

func (n *node) addEntry(e *Entry) {
	n.entries = append(n.entries, e)
}

func (t *CodeOwners) FindOwners(path string) []string {
	owners := []string{}
	walker := func(key string, value interface{}) error {
		if value == nil {
			return nil
		}
		n, ok := value.(*node)
		if !ok {
			panic("Structure of the code owner index is malformed")
		}

		for _, en := range n.entries {
			if en.suffix == PathSufix(Recursive) || en.suffix == PathSufix(Absolute) {
				owners = append(owners, en.owners...)
			}
		}

		return nil
	}
	t.WalkKey(path, walker)

	//get the base wild card type
	ext := "*" + filepath.Ext(path)
	extEntry := t.Get(ext)
	if extEntry != nil {
		n, ok := extEntry.(*node)
		if !ok {
			log.Fatal("Structure of the code owner index is malformed")
		}
		for _, en := range n.entries {
			owners = append(owners, en.owners...)
		}
	}
	return removeDuplicatesUnordered(owners)
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}
