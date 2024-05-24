package smanchai

import (
	"strconv"
)

type Visitor struct {
	emitter chan *emitted
	constc  int
	state   int
}

func NewVisitor() *Visitor {
	return &Visitor{
		emitter: make(chan *emitted),
		constc:  0,
	}
}

func (v *Visitor) Emitter() chan *emitted {
	return v.emitter
}

func consInst(opcode int) *emitted {
	return &emitted{
		Type:  E_Inst,
		Value: opcode,
	}
}

func (v *Visitor) Accept(node *Node) {
	if node == nil {
		v.emitter <- nil
		return
	}
	switch node.Type {
	case AstProgram:
		v.visitProgram(node.Object.(*ProgramNode))
	case AstExpression:
		v.visitExpression(node.Object.(*ExpressionNode))
	case AstDisjunction:
		v.visitDisjunction(node.Object.(*DisjunctionNode))
	case AstConjunction:
		v.visitConjunction(node.Object.(*ConjunctionNode))
	case AstEquality:
		v.visitEquality(node.Object.(*ComparisonNode))
	case AstComparison:
		v.visitComparison(node.Object.(*ComparisonNode))
	case AstPrimitive:
		v.visitPrimitive(node)
	case AstIdentifier:
		v.visitIdentifier(node.Object.(*IdentifierNode))
	case AstLiteral:
		v.visitLiteral(node.Object.(*LiteralNode))
	default:
		panic("I forgot to implement Visitor.")
	}
}

func (v *Visitor) visitProgram(node *ProgramNode) {
	for i := 0; i < len(node.Children); i++ {
		v.Accept(node.Children[i])
	}
	v.emitter <- nil
}

func (v *Visitor) visitExpression(node *ExpressionNode) {
	m := false
	v.Accept(node.Left)
	if v.state == DTypeString {
		m = true
	}
	v.state = -1
	v.Accept(node.Right)
	if v.state == DTypeString {
		m = true
	}
	v.state = -1
	switch node.Op {
	case EOpADD:
		if m {
			v.emitter <- consInst(Op_sconcat)
		} else {
			v.emitter <- consInst(Op_iadd)
		}
	case EOprSUB:
		v.emitter <- consInst(Op_dsub)
	case EOprMULT:
		v.emitter <- consInst(Op_dmul)
	case EOprDIV:
		v.emitter <- consInst(Op_ddiv)
	case EOprPOW:
		v.emitter <- consInst(Op_dexp)
	}
}

func (v *Visitor) visitDisjunction(node *DisjunctionNode) {
	v.Accept(node.Left)
	v.Accept(node.Right)
	v.emitter <- consInst(Op_ior)
}

func (v *Visitor) visitConjunction(node *ConjunctionNode) {
	v.Accept(node.Left)
	v.Accept(node.Right)
	v.emitter <- consInst(Op_iand)
}

func (v *Visitor) visitEquality(node *ComparisonNode) {
	v.Accept(node.Left)
	v.Accept(node.Right)
	switch node.Op {
	case OprEqual:
		v.emitter <- consInst(Op_cmp_eq)
	case OprNotEqual:
		v.emitter <- consInst(Op_cmp_ne)
	}
}

func (v *Visitor) visitComparison(node *ComparisonNode) {
	v.Accept(node.Left)
	v.Accept(node.Right)
	switch node.Op {
	case OprGreaterThan:
		v.emitter <- consInst(Op_cmp_g)
	case OprGreaterThanEqual:
		v.emitter <- consInst(Op_cmp_ge)
	case OprLessThan:
		v.emitter <- consInst(Op_cmp_l)
	case OprLessThanEqual:
		v.emitter <- consInst(Op_cmp_le)
	}
}

func (v *Visitor) visitPrimitive(node *Node) {
	v.Accept(node.Object.(*Node))
}

func (v *Visitor) visitIdentifier(node *IdentifierNode) {
	if node.At {
		Const := &emitted{Type: E_Const}
		ConstData := &Data{Type: DTypeDataRef, Value: &DataRef{
			Name: node.Base,
			Root: DRTypeVMStatic,
		}}
		Const.Value = ConstData
		v.emitter <- Const
		v.emitter <- consInst(Op_getstatic)
		v.emitter <- consInst(v.constc)
		v.constc++
		for _, attr := range node.SubIdentifier {
			Const := &emitted{Type: E_Const}
			ConstData := &Data{Type: DTypeAttrRef, Value: &AttrRef{
				Name: attr,
				Root: v.constc - 1,
			}}
			Const.Value = ConstData
			v.emitter <- Const
			v.emitter <- consInst(Op_getattr)
			v.emitter <- consInst(v.constc)
			v.constc++
		}
	}
}

func (v *Visitor) visitLiteral(node *LiteralNode) {
	Const := &emitted{Type: E_Const}
	ConstData := &Data{Value: node.Raw}
	Const.Value = ConstData
	switch node.Type {
	case LSTRING:
		ConstData.Type = DTypeString
		v.emitter <- Const
		v.emitter <- consInst(Op_sload)
		v.emitter <- consInst(v.constc)
	case LBOOLEAN:
		ConstData.Type = DTypeBool
		f := 0
		if node.Raw == "true" {
			f = 1
		}
		ConstData.Value = f
		v.emitter <- Const
		v.emitter <- consInst(Op_iload)
		v.emitter <- consInst(v.constc)
	case LNUMBER:
		ConstData.Type = DTypeDouble
		f, err := strconv.ParseFloat(node.Raw, 64)
		if err != nil {
			panic(err)
		}
		ConstData.Value = f
		v.emitter <- Const
		v.emitter <- consInst(Op_dload)
		v.emitter <- consInst(v.constc)
	}
	v.constc++
}

func Compile(node *Node) *VM {
	consts := make([]*Data, 0, 256)
	insts := make([]int, 0, 256)
	visitor := NewVisitor()
	emit := visitor.Emitter()
	go visitor.Accept(node)
	for {
		obj := <-emit
		if obj == nil {
			break
		}
		if obj.Type == E_Inst {
			insts = append(insts, obj.Value.(int))
		}
		if obj.Type == E_Const {
			visitor.state = int(obj.Value.(*Data).Type)
			consts = append(consts, obj.Value.(*Data))
		}
	}
	return &VM{
		ConstsPool: consts,
		Insts:      insts,
		pc:         0,
		static:     map[string]static{},
	}
}
