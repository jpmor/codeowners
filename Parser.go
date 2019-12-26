package codeowners

import (
	"bytes"
	"io"
)

type Entry struct {
	path    string
	suffix  int
	comment string
	owners  []string
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

func newEntry() *Entry {
	return &Entry{
		owners: make([]string, 0),
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a line from a codeowners file.
func (p *Parser) Parse() (*Entry, error) {
	entry := newEntry()

	tok, lit := p.scanIgnoreWhitespace()

	//Capture the comment (entire line is a comment)
	if tok == HASH {
		var b bytes.Buffer
		b.WriteString(lit)
		for {
			tok, lit = p.scan()
			if tok == EOF {
				break
			}
			b.WriteString(lit)
		}
		entry.comment = b.String()
		return entry, nil
	}

	var b bytes.Buffer
	for tok != WS {
		b.WriteString(lit)
		tok, lit = p.scan()
	}
	//TODO: Validate/normalize file path?
	entry.path = b.String()

	tok, lit = p.scanIgnoreWhitespace()
	for tok != EOF {
		b.Reset()
		for tok != WS {
			if tok == EOF {
				break
			}
			b.WriteString(lit)
			tok, lit = p.scan()
		}
		entry.owners = append(entry.owners, b.String())
		tok, lit = p.scanIgnoreWhitespace()
	}

	// Return the successfully parsed statement.
	return entry, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
