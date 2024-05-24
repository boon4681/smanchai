package smanchai

type Token int

const (
	EOF                 = iota // ✅
	ILLEGAL                    // ✅
	WS                         // ✅
	IDENTIFIER                 // ✅
	STRING                     // ✅
	NUMBER                     // ✅
	BOOL                       // ✅
	AT                         // ✅
	ADD                        // ✅
	SUB                        // ✅
	MULT                       // ✅
	DIV                        // ✅
	POW                        // ✅
	LParent                    // ✅
	RParent                    // ✅
	COMMA                      // ✅
	DOT                        // ✅
	EQUALITY_OPERATOR          // ✅
	COMPARISON_OPERATOR        // ✅
	CONJUNCTION                // ✅
	DISJUNCTION                // ✅
)

var tokens = []string{
	EOF:        "EOF",
	ILLEGAL:    "ILLEGAL",
	WS:         "WS",
	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	NUMBER:     "NUMBER",
	BOOL:       "BOOL",
	AT:         "AT",
	ADD:        "ADD",
	SUB:        "SUB",
	MULT:       "MULT",
	DIV:        "DIV",
	POW:        "POW",
	LParent:    "LParent",
	RParent:    "RParent",
	COMMA:      "COMMA",
	DOT:        "DOT",

	EQUALITY_OPERATOR:   "EQUALITY_OPERATOR",
	COMPARISON_OPERATOR: "COMPARISON_OPERATOR",

	CONJUNCTION: "CONJUNCTION",
	DISJUNCTION: "DISJUNCTION",
}

func (t Token) String() string {
	return tokens[t]
}

type TokenInfo struct {
	r     Range
	token Token
	str   string
}
