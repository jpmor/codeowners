package codeowners

type Token int

const (
	ILLEGAL Token = iota
	EOF
	EOL
	WS

	// Literals
	IDENT // main

	// Misc characters
	ASTERISK // *
	COMMA    // ,
	HASH     // #
)
