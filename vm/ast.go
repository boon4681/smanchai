package smanchai

type AstType int

const (
	AstProgram = iota
	AstIdentifier
	AstPrimitive
	AstLiteral
	AstEquality
	AstComparison
	AstExpression
	AstConjunction
	AstDisjunction
)

var ast = []string{
	AstProgram:     "AstProgram",
	AstIdentifier:  "AstIdentifier",
	AstPrimitive:   "AstPrimitive",
	AstLiteral:     "AstLiteral",
	AstEquality:    "AstEquality",
	AstComparison:  "AstComparison",
	AstExpression:  "AstExpression",
	AstConjunction: "AstConjunction",
	AstDisjunction: "AstDisjunction",
}

func (s AstType) String() string {
	return ast[s]
}

type LiteralType int

const (
	LSTRING = iota
	LNUMBER
	LBOOLEAN
)

type EquationOpr int

const (
	EOpADD = iota
	EOprSUB
	EOprMULT
	EOprDIV
	EOprPOW
)

type ComparisonOpr int

const (
	OprLessThan = iota
	OprGreaterThan
	OprLessThanEqual
	OprGreaterThanEqual
	OprEqual
	OprNotEqual
)

type Node struct {
	Type   AstType
	Range  Range
	Object any
}

type ProgramNode struct {
	Children []*Node
}

type IdentifierNode struct {
	At            bool
	Base          string
	SubIdentifier []string
}

type LiteralNode struct {
	Type LiteralType
	Raw  string
}

type FunctionNode struct {
	Name   string
	Params []*Node
}

type ComparisonNode struct {
	Left  *Node
	Op    ComparisonOpr
	Right *Node
}

type ExpressionNode struct {
	Left  *Node
	Op    EquationOpr
	Right *Node
}

type DisjunctionNode struct {
	Left  *Node
	Right *Node
}

type ConjunctionNode struct {
	Left  *Node
	Right *Node
}
