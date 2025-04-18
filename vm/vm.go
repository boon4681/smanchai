package smanchai

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type Opcode int

const (
	Op = iota
	Op_getstatic
	Op_getattr
	Op_invoke    // invoke built-in or extended function
	Op_loadrange // load Token range for good error
	Op_iload     // load int from const to stack
	Op_iinc      // increase int at index <x> by <y>
	Op_iadd      //
	Op_isub      //
	Op_imul      //
	Op_idiv      //
	Op_imod      //
	Op_iand      //
	Op_ior       //
	Op_i2b       // int to boolean
	Op_i2c       // int to char
	Op_i2d       // int to double

	Op_dload //
	Op_dinc  //
	Op_dadd  //
	Op_dsub  //
	Op_dmul  //
	Op_ddiv  //
	Op_dmod  //
	Op_dexp  //

	Op_sload   // load string from const to stack
	Op_sconcat // concat string
	Op_cmp_eq  // compare eq
	Op_cmp_ne  // compare not eq
	Op_cmp_g   // compare greater
	Op_cmp_ge  // compare greater eq
	Op_cmp_l   // compare less
	Op_cmp_le  // compare less eq
)

var opcodes = []string{
	Op_getstatic: "Op_getstatic",
	Op_getattr:   "Op_getattr",
	Op_invoke:    "Op_invoke",
	Op_iload:     "Op_iload",
	Op_iinc:      "Op_iinc",
	Op_iadd:      "Op_iadd",
	Op_isub:      "Op_isub",
	Op_imul:      "Op_imul",
	Op_idiv:      "Op_idiv",
	Op_imod:      "Op_imod",
	Op_iand:      "Op_iand",
	Op_ior:       "Op_ior",
	Op_i2b:       "Op_i2b",
	Op_i2c:       "Op_i2c",
	Op_i2d:       "Op_i2d",
	Op_dload:     "Op_dload",
	Op_dinc:      "Op_dinc",
	Op_dadd:      "Op_dadd",
	Op_dsub:      "Op_dsub",
	Op_dmul:      "Op_dmul",
	Op_ddiv:      "Op_ddiv",
	Op_dmod:      "Op_dmod",
	Op_dexp:      "Op_dexp",
	Op_sload:     "Op_sload",
	Op_sconcat:   "Op_sconcat",
	Op_cmp_eq:    "Op_cmp_eq",
	Op_cmp_ne:    "Op_cmp_ne",
	Op_cmp_g:     "Op_cmp_g",
	Op_cmp_ge:    "Op_cmp_ge",
	Op_cmp_l:     "Op_cmp_l",
	Op_cmp_le:    "Op_cmp_le",
}

func (op Opcode) String() string {
	return opcodes[op]
}

type emittedType int

const (
	E_Inst = iota
	E_Const
)

type emitted struct {
	Type  emittedType
	Value any
}

type dataType int

const (
	DTypeInst = iota
	DTypeArray
	DTypeInt
	DTypeDouble
	DTypeBool
	DTypeChar
	DTypeString
	DTypeDataRef
	DTypeAttrRef
	DTypeMethodRef
	DTypeObject
)

type DataRefType int

const (
	DRTypeVMStatic = iota
	DRTypeGlobal
	DRTypeLocal
)

type DataRef struct {
	Name string
	Root DataRefType
}

type AttrRef struct {
	Name string
	Root int
}

type Data struct {
	Type  dataType
	Value any
}

type DataObjectArray struct {
	Data []*Data
}

type DataObjectMap struct {
	Data map[string]*Data
}

func (t *Data) String() string {
	return fmt.Sprintf("%s", t.Value)
}

type stack[T any] struct {
	arr []T
}

func (s *stack[T]) Push(v T) {
	s.arr = append(s.arr, v)
}

func (s *stack[T]) Pop() T {
	l := len(s.arr)
	if l == 0 {
		panic("stack no item left to pop")
	}
	v := s.arr[l-1]
	s.arr = s.arr[:l-1]
	return v
}

func (s *stack[T]) Get(i int) T {
	l := len(s.arr) - 1
	return s.arr[l-i]
}

func (s *stack[T]) Len() int {
	return len(s.arr)
}

type static func() *Data

type VM struct {
	static     map[string]static
	ConstsPool []*Data
	Insts      []int
	pc         int // program counter
	Dg         bool
}

func (vm *VM) AddStatic(name string, ex static) {
	vm.static[name] = ex
}

func (vm *VM) Run() (*Data, error) {
	operand := stack[*Data]{arr: make([]*Data, 0, 256)}
	if vm.Dg {
		fmt.Println("run vm")
	}
	for {
		if vm.pc >= len(vm.Insts) {
			break
		}
		inst := vm.Insts[vm.pc]
		if vm.Dg {
			fmt.Printf("Rx:\t%d:%d\t%s\n", vm.pc, inst, opcodes[inst])
			{
				b, err := json.Marshal(operand.arr)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(string(b))
				}
			}
		}
		switch inst {
		case Op_getstatic:
			vm.pc++
			b0 := vm.ConstsPool[vm.Insts[vm.pc]]
			switch b0.Type {
			case DTypeDataRef:
				b1 := b0.Value.(*DataRef)
				switch b1.Root {
				case DRTypeVMStatic:
					if v, o := vm.static[b1.Name]; o {
						operand.Push(v())
					} else {
						return nil, fmt.Errorf("error: not found static named \"%s\"", b1.Name)
					}
				default:
					return nil, fmt.Errorf("unknown error: stack error")
				}
			default:
				return nil, fmt.Errorf("unknown error: stack error")
			}
		case Op_getattr:
			vm.pc++
			b0 := vm.ConstsPool[vm.Insts[vm.pc]]
			switch b0.Type {
			case DTypeAttrRef:
				b1 := b0.Value.(*AttrRef)
				b2 := operand.Pop()
				switch b2.Type {
				case DTypeArray:
					i, err := strconv.Atoi(b1.Name)
					if err != nil {
						return nil, err
					}
					b3 := b2.Value.(*DataObjectArray)
					operand.Push(b3.Data[i])
				case DTypeObject:
					b3 := b2.Value.(*DataObjectMap)
					operand.Push(b3.Data[b1.Name])
				default:
					return nil, fmt.Errorf("unknown error: wrong data type of getattr operand")
				}
			default:
				return nil, fmt.Errorf("unknown error: stack error")
			}
		case Op_invoke:
			vm.pc++
			b0 := vm.ConstsPool[vm.Insts[vm.pc]]
			switch b0.Type {
			case DTypeMethodRef:
			default:
				return nil, fmt.Errorf("unknown error: stack error")
			}
		case Op_iload:
			vm.pc++
			b0 := vm.ConstsPool[vm.Insts[vm.pc]]
			if b0.Type == DTypeInt || b0.Type == DTypeBool || b0.Type == DTypeChar {
				operand.Push(b0)
			} else {
				return nil, fmt.Errorf("unknown error: stack error")
			}
		case Op_iinc:
			vm.pc++
			a0 := vm.Insts[vm.pc]
			a1 := vm.Insts[vm.pc]
			b0 := vm.ConstsPool[a0]
			if b0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			b0.Value = b0.Value.(int) + a1
		case Op_iadd, Op_isub, Op_imul, Op_idiv, Op_imod:
			a0 := operand.Pop()
			a1 := operand.Pop()
			if a0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			if a1.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			switch inst {
			case Op_iadd:
				a0.Value = a1.Value.(int) + a0.Value.(int)
			case Op_isub:
				a0.Value = a1.Value.(int) - a0.Value.(int)
			case Op_imul:
				a0.Value = a1.Value.(int) * a0.Value.(int)
			case Op_idiv:
				a0.Value = a1.Value.(int) / a0.Value.(int)
			case Op_imod:
				a0.Value = a1.Value.(int) % a0.Value.(int)
			}
			operand.Push(a0)
		case Op_iand, Op_ior:
			a0 := operand.Pop()
			a1 := operand.Pop()
			if a0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			if a1.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			switch inst {
			case Op_iand:
				operand.Push(&Data{
					Type:  DTypeInt,
					Value: a0.Value.(int) & a1.Value.(int),
				})
			case Op_ior:
				operand.Push(&Data{
					Type:  DTypeInt,
					Value: a0.Value.(int) | a1.Value.(int),
				})
			}
			// TODO
		case Op_i2b:
			a0 := operand.Pop()
			if a0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			m := a0.Value.(int)
			if m > 0 {
				m = 1
			} else {
				m = 0
			}
			operand.Push(&Data{
				Type:  DTypeBool,
				Value: m,
			})
		case Op_i2c:
			a0 := operand.Pop()
			if a0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			operand.Push(&Data{
				Type:  DTypeChar,
				Value: byte(a0.Value.(int)),
			})
		case Op_i2d:
			a0 := operand.Pop()
			if a0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			operand.Push(&Data{
				Type:  DTypeChar,
				Value: float64(a0.Value.(int)),
			})
		case Op_dload:
			vm.pc++
			b0 := vm.ConstsPool[vm.Insts[vm.pc]]
			if b0.Type == DTypeDouble {
				operand.Push(b0)
			} else {
				return nil, fmt.Errorf("unknown error: stack error")
			}
		case Op_dinc:
			vm.pc++
			a0 := vm.Insts[vm.pc]
			a1 := vm.Insts[vm.pc]
			b0 := vm.ConstsPool[a0]
			if b0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			b0.Value = b0.Value.(float64) + float64(a1)
		case Op_dadd, Op_dsub, Op_dmul, Op_ddiv, Op_dmod, Op_dexp:
			a0 := operand.Pop()
			a1 := operand.Pop()
			if a0.Type != DTypeDouble && a0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			if a1.Type != DTypeDouble && a0.Type != DTypeInt {
				return nil, fmt.Errorf("unknown error: stack error")
			}
			switch inst {
			case Op_dadd:
				a0.Value = a1.Value.(float64) + a0.Value.(float64)
			case Op_dsub:
				a0.Value = a1.Value.(float64) - a0.Value.(float64)
			case Op_dmul:
				a0.Value = a1.Value.(float64) * a0.Value.(float64)
			case Op_ddiv:
				a0.Value = a1.Value.(float64) / a0.Value.(float64)
			case Op_dmod:
				a0.Value = a1.Value.(int) % a0.Value.(int)
			case Op_dexp:
				a0.Value = math.Pow(a1.Value.(float64), a0.Value.(float64))
			}
			operand.Push(a0)
		case Op_sload:
			vm.pc++
			b0 := vm.ConstsPool[vm.Insts[vm.pc]]
			if b0.Type == DTypeString {
				operand.Push(b0)
			} else {
				return nil, fmt.Errorf("unknown error: string")
			}
		case Op_sconcat:
			a0 := operand.Pop()
			a1 := operand.Pop()
			if a0 != nil && a0.Type != DTypeString {
				return nil, fmt.Errorf("unknown error: you cannot concat string with other type")
			}
			if a1 != nil && a1.Type != DTypeString {
				return nil, fmt.Errorf("unknown error: you cannot concat string with other type")
			}
			var a, b string
			if a0 == nil {
				a = "undefined"
			} else {
				a = a0.Value.(string)
			}
			if a1 == nil {
				b = "undefined"
			} else {
				b = a1.Value.(string)
			}
			a0.Value = b + a
			operand.Push(a0)
		case Op_cmp_eq, Op_cmp_ne, Op_cmp_g, Op_cmp_ge, Op_cmp_l, Op_cmp_le:
			a0 := operand.Pop()
			a1 := operand.Pop()
			var v int = 0
			switch inst {
			case Op_cmp_eq:
				if a0.Value == a1.Value {
					v = 1
				}
			case Op_cmp_ne:
				if a0.Value != a1.Value {
					v = 1
				}
			case Op_cmp_g:
				if (a0.Type == DTypeInt || a0.Type == DTypeChar || a0.Type == DTypeDouble || a0.Type == DTypeBool) &&
					(a1.Type == DTypeInt || a1.Type == DTypeChar || a1.Type == DTypeDouble || a1.Type == DTypeBool) {
					if toFloat64(a1.Value) > toFloat64(a0.Value) {
						v = 1
					}
				}
			case Op_cmp_ge:
				if (a0.Type == DTypeInt || a0.Type == DTypeChar || a0.Type == DTypeDouble || a0.Type == DTypeBool) &&
					(a1.Type == DTypeInt || a1.Type == DTypeChar || a1.Type == DTypeDouble || a1.Type == DTypeBool) {
					if toFloat64(a1.Value) >= toFloat64(a0.Value) {
						v = 1
					}
				}
			case Op_cmp_l:
				if (a0.Type == DTypeInt || a0.Type == DTypeChar || a0.Type == DTypeDouble || a0.Type == DTypeBool) &&
					(a1.Type == DTypeInt || a1.Type == DTypeChar || a1.Type == DTypeDouble || a1.Type == DTypeBool) {
					if toFloat64(a1.Value) < toFloat64(a0.Value) {
						v = 1
					}
				}
			case Op_cmp_le:
				if (a0.Type == DTypeInt || a0.Type == DTypeChar || a0.Type == DTypeDouble || a0.Type == DTypeBool) &&
					(a1.Type == DTypeInt || a1.Type == DTypeChar || a1.Type == DTypeDouble || a1.Type == DTypeBool) {
					if toFloat64(a1.Value) <= toFloat64(a0.Value) {
						v = 1
					}
				}
			}
			operand.Push(&Data{
				Type:  DTypeInt,
				Value: v,
			})
		}
		vm.pc++
	}
	if operand.Len() > 0 {
		result := operand.Pop()
		if vm.Dg {
			fmt.Printf("Result %s\n", result)
		}
		return result, nil
	}
	return nil, nil
}

func toFloat64(value interface{}) float64 {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint())
	default:
		return 0
	}
}

func TestVM() *VM {
	consts := make([]*Data, 0)
	insts := make([]int, 0)
	consts = append(consts, &Data{
		Type:  DTypeInt,
		Value: 8,
	})
	consts = append(consts, &Data{
		Type:  DTypeInt,
		Value: 10,
	})
	consts = append(consts, &Data{
		Type:  DTypeInt,
		Value: 7,
	})
	insts = append(insts, Op_iload)
	insts = append(insts, 2)
	insts = append(insts, Op_iload)
	insts = append(insts, 0)
	insts = append(insts, Op_iload)
	insts = append(insts, 1)
	insts = append(insts, Op_imul)
	insts = append(insts, Op_iadd)
	return &VM{
		ConstsPool: consts,
		Insts:      insts,
		pc:         0,
	}
}
