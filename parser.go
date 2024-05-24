package smanchai

import "fmt"

type Parser struct {
	lexer     *Lexer
	current   *TokenInfo
	prev      *TokenInfo
	on_unnext bool
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer:   lexer,
		current: nil,
		prev:    nil,
	}
}

func (p *Parser) next() (Range, Token, string) {
	if p.on_unnext {
		current := *p.current
		p.prev = &current
		p.current = nil
		p.on_unnext = false
		return current.r, current.token, current.str
	} else {
		r, token, str := p.lexer.Lex()
		if token == ILLEGAL {
			panic(fmt.Errorf("Unknown syntax: at %s \"%s\"", r.String(), str))
		}
		p.current = &TokenInfo{
			r:     r,
			token: token,
			str:   str,
		}
		prev := *p.current
		p.prev = &prev
		return r, token, str
	}
}

func (p *Parser) skip_whitespace() {
	_, token, _ := p.next()
	for ; token == WS && token != EOF; _, token, _ = p.next() {
	}
	p.unnext()
}

func (p *Parser) unnext() {
	if p.prev == nil && p.current == nil {
		panic("Unnext cannot be used if the last called was not next method")
	}
	p.on_unnext = true
	prev := *p.prev
	p.current = &prev
	p.prev = nil
}

func (p *Parser) Parse() *Node {
	_, token, _ := p.next()
	if token == EOF {
		return nil
	}
	p.unnext()
	return p.astProgram()
}

func (p *Parser) astProgram() *Node {
	obj := &Node{
		Type: AstProgram,
	}
	program := &ProgramNode{}
	obj.Object = program
	children := make([]*Node, 0, 1)
	p.skip_whitespace()
	if o := p.astDisjunction(); o != nil {
		children = append(children, o)
		program.Children = children
		return obj
	}
	return nil
}

func (p *Parser) astDisjunction() *Node {
	obj := &Node{
		Type: AstDisjunction,
	}
	disj := &DisjunctionNode{}
	obj.Object = disj
	if left := p.astConjunction(); left != nil {
		p.skip_whitespace()
		r, token, _ := p.next()
		obj.Range = r
		switch token {
		case DISJUNCTION:
			p.skip_whitespace()
			if right := p.astConjunction(); right != nil {
				disj.Left = left
				disj.Right = right
				return obj
			} else {
				panic("")
			}
		default:
			p.unnext()
			return left
		}
	}
	return nil
}

func (p *Parser) astConjunction() *Node {
	obj := &Node{
		Type: AstConjunction,
	}
	conj := &ConjunctionNode{}
	obj.Object = conj
	if left := p.astEqaulity(); left != nil {
		p.skip_whitespace()
		r, token, _ := p.next()
		obj.Range = r
		switch token {
		case CONJUNCTION:
			p.skip_whitespace()
			if right := p.astEqaulity(); right != nil {
				conj.Left = left
				conj.Right = right
				return obj
			} else {
				panic("")
			}
		default:
			p.unnext()
			return left
		}
	}
	return nil
}

func (p *Parser) astEqaulity() *Node {
	obj := &Node{
		Type: AstEquality,
	}
	expr := &ComparisonNode{}
	obj.Object = expr
	if left := p.astComparison(); left != nil {
		p.skip_whitespace()
		r, token, str := p.next()
		obj.Range = r
		expr.Left = left
		switch token {
		case EQUALITY_OPERATOR:
			switch str {
			case "==":
				expr.Op = OprEqual
			case "!=":
				expr.Op = OprNotEqual
			default:
				p.unnext()
				return left
			}
			p.skip_whitespace()
			if right := p.astComparison(); right != nil {
				expr.Right = right
				return obj
			} else {
				panic("")
			}
		default:
			p.unnext()
			return left
		}
	}
	return nil
}

func (p *Parser) astComparison() *Node {
	obj := &Node{
		Type: AstComparison,
	}
	expr := &ComparisonNode{}
	obj.Object = expr
	if left := p.astAdditiveExpression(); left != nil {
		p.skip_whitespace()
		r, token, str := p.next()
		obj.Range = r
		expr.Left = left
		switch token {
		case COMPARISON_OPERATOR:
			switch str {
			case ">":
				expr.Op = OprGreaterThan
			case "<":
				expr.Op = OprLessThan
			case ">=":
				expr.Op = OprGreaterThanEqual
			case "<=":
				expr.Op = OprLessThanEqual
			default:
				p.unnext()
				return left
			}
			p.skip_whitespace()
			if right := p.astAdditiveExpression(); right != nil {
				expr.Right = right
				return obj
			} else {
				panic("")
			}
		default:
			p.unnext()
		}
		return left
	}
	return nil
}

func (p *Parser) astAdditiveExpression() *Node {
	obj := &Node{
		Type: AstExpression,
	}
	expr := &ExpressionNode{}
	obj.Object = expr
	if left := p.astMultiplicativeExpression(); left != nil {
		p.skip_whitespace()
		r, token, _ := p.next()
		obj.Range = r
		expr.Left = left
		switch token {
		case ADD:
			expr.Op = EOpADD
			p.skip_whitespace()
			if right := p.astMultiplicativeExpression(); right != nil {
				expr.Right = right
				return obj
			} else {
				panic("")
			}
		case SUB:
			expr.Op = EOprSUB
			p.skip_whitespace()
			if right := p.astMultiplicativeExpression(); right != nil {
				expr.Right = right
				return obj
			} else {
				panic("")
			}
		default:
			p.unnext()
		}
		return left
	}
	return nil
}

func (p *Parser) astMultiplicativeExpression() *Node {
	obj := &Node{
		Type: AstExpression,
	}
	expr := &ExpressionNode{}
	obj.Object = expr
	if left := p.astExponentialExpression(); left != nil {
		p.skip_whitespace()
		r, token, _ := p.next()
		obj.Range = r
		expr.Left = left
		switch token {
		case MULT:
			expr.Op = EOprMULT
			p.skip_whitespace()
			if right := p.astExponentialExpression(); right != nil {
				expr.Right = right
				return obj
			} else {
				panic("")
			}
		case DIV:
			expr.Op = EOprDIV
			p.skip_whitespace()
			if right := p.astExponentialExpression(); right != nil {
				expr.Right = right
				return obj
			} else {
				panic("")
			}
		default:
			p.unnext()
		}
		return left
	}
	return nil
}

func (p *Parser) astExponentialExpression() *Node {
	obj := &Node{
		Type: AstExpression,
	}
	expr := &ExpressionNode{}
	obj.Object = expr
	if left := p.astPrimitive(); left != nil {
		switch left.Type {
		case AstPrimitive:
			p.skip_whitespace()
			if r, token, _ := p.next(); token == POW {
				p.skip_whitespace()
				expr.Left = left
				if right := p.astPrimitive(); right != nil {
					obj.Range = r
					expr.Op = EOprPOW
					expr.Right = right
					return obj
				} else {
					panic("")
				}
			} else {
				p.unnext()
				return left
			}
		case AstExpression:
			return left
		}
	}
	return nil
}

func (p *Parser) astPrimitive() *Node {
	p.skip_whitespace()
	r, token, _ := p.next()
	if token == EOF {
		return nil
	}
	p.unnext()
	obj := &Node{
		Type:  AstPrimitive,
		Range: r,
	}
	if o := p.astFunction(); o != nil {
		obj.Object = o
		return obj
	}
	if o := p.astIdentifier(); o != nil {
		obj.Object = o
		return obj
	}
	if o := p.astLiteral(); o != nil {
		obj.Object = o
		return obj
	}
	return nil
}

func (p *Parser) astLiteral() *Node {
	r, token, str := p.next()
	obj := &Node{
		Type:  AstLiteral,
		Range: r,
	}
	if token == NUMBER {
		obj.Object = &LiteralNode{
			Type: LNUMBER,
			Raw:  str,
		}
		return obj
	}
	if token == STRING {
		obj.Object = &LiteralNode{
			Type: LSTRING,
			Raw:  str,
		}
		return obj
	}
	if token == BOOL {
		obj.Object = &LiteralNode{
			Type: LBOOLEAN,
			Raw:  str,
		}
		return obj
	}
	return nil
}

func (p *Parser) astFunction() *Node {
	return nil
}

func (p *Parser) astIdentifier() *Node {
	r, token, str := p.next()
	base := ""
	subIdentifier := make([]string, 0, 256)
	if token == AT {
		if r, token, str = p.next(); token == IDENTIFIER {
			base = str
		}
		for {
			if _, token, _ := p.next(); token != DOT {
				p.unnext()
				break
			}
			if r, token, str = p.next(); token != IDENTIFIER {
				if token == WS {
					if r, token, str = p.next(); token == IDENTIFIER {
						subIdentifier = append(subIdentifier, str)
					} else {
						panic(fmt.Sprintf("Error: Invalid or unexpected token at %d:%d, \"%s\"\n", r.Line, r.Column, str))
					}
				} else {
					break
				}
			} else {
				subIdentifier = append(subIdentifier, str)
			}
		}
		return &Node{
			Type:  AstIdentifier,
			Range: r,
			Object: &IdentifierNode{
				At:            true,
				Base:          base,
				SubIdentifier: subIdentifier,
			},
		}
	}
	if token == IDENTIFIER {
		base = str
		for {
			if _, token, _ := p.next(); token != DOT {
				p.unnext()
				break
			}
			if r, token, str = p.next(); token != IDENTIFIER {
				if token == WS {
					if r, token, str = p.next(); token == IDENTIFIER {
						subIdentifier = append(subIdentifier, str)
					} else {
						panic(fmt.Sprintf("Error: Invalid or unexpected token at %d:%d, \"%s\"\n", r.Line, r.Column, str))
					}
				} else {
					break
				}
			} else {
				subIdentifier = append(subIdentifier, str)
			}
		}
		return &Node{
			Type:  AstIdentifier,
			Range: r,
			Object: &IdentifierNode{
				At:            false,
				Base:          base,
				SubIdentifier: subIdentifier,
			},
		}
	}
	p.unnext()
	return nil
}
