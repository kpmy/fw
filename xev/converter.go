package xev

import (
	"fmt"
	"fw/cp/constant"
	"fw/cp/constant/enter"
	"fw/cp/constant/operation"
	"fw/cp/constant/statement"
	"fw/cp/module"
	"fw/cp/node"
	"fw/cp/object"
	"math/big"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
	"ypk/assert"
)

func (r *Result) findNode(id string) *Node {
	var ret *Node
	for i := 0; i < len(r.GraphList) && (ret == nil); i++ {
		for j := 0; j < len(r.GraphList[i].NodeList) && (ret == nil); j++ {
			if id == r.GraphList[i].NodeList[j].Id {
				ret = &r.GraphList[i].NodeList[j]
			}
		}
	}
	return ret
}

func (r *Result) findLink(n *Node, link string) *Node {
	target := ""
	for i := range r.GraphList {
		for j := range r.GraphList[i].EdgeList {
			if (r.GraphList[i].EdgeList[j].Source == n.Id) && (r.GraphList[i].EdgeList[j].CptLink == link) {
				target = r.GraphList[i].EdgeList[j].Target
			}
		}
	}
	var ret *Node
	if target != "" {
		ret = r.findNode(target)
	}
	return ret
}

type convertable interface {
	SetData(x interface{})
	SetType(t object.Type)
}

type typed interface {
	SetType(t object.Type)
}

func initType(typ string, conv typed) {
	assert.For(conv != nil, 20)
	switch typ {
	case "INTEGER":
		conv.SetType(object.INTEGER)
	case "SHORTINT":
		conv.SetType(object.SHORTINT)
	case "LONGINT":
		conv.SetType(object.LONGINT)
	case "BYTE":
		conv.SetType(object.BYTE)
	case "CHAR":
		conv.SetType(object.CHAR)
	case "SHORTCHAR":
		conv.SetType(object.SHORTCHAR)
	case "REAL":
		conv.SetType(object.REAL)
	case "SHORTREAL":
		conv.SetType(object.SHORTREAL)
	case "SET":
		conv.SetType(object.SET)
	case "BOOLEAN":
		conv.SetType(object.BOOLEAN)
	case "":
		conv.SetType(object.NOTYPE)
	default:
		panic("no such constant type")
	}
}

func convertData(typ string, val string, conv convertable) {
	assert.For(conv != nil, 20)
	switch typ {
	case "INTEGER":
		conv.SetType(object.INTEGER)
		x, _ := strconv.ParseInt(val, 10, 32)
		conv.SetData(int32(x))
	case "SHORTINT":
		conv.SetType(object.SHORTINT)
		x, _ := strconv.ParseInt(val, 10, 16)
		conv.SetData(int16(x))
	case "LONGINT":
		conv.SetType(object.LONGINT)
		x, _ := strconv.ParseInt(val, 10, 64)
		conv.SetData(x)
	case "BYTE":
		conv.SetType(object.BYTE)
		x, _ := strconv.ParseInt(val, 10, 8)
		conv.SetData(int8(x))
	case "CHAR":
		conv.SetType(object.CHAR)
		x, _ := strconv.ParseUint(val, 16, 16)
		s := make([]uint16, 1)
		s[0] = uint16(x)
		r := utf16.Decode(s)
		conv.SetData(r[0])
	case "SHORTCHAR":
		conv.SetType(object.SHORTCHAR)
		x, _ := strconv.ParseUint(val, 16, 8)
		s := make([]uint8, 1)
		s[0] = uint8(x)
		r, _ := utf8.DecodeRune(s)
		conv.SetData(r)
	case "REAL":
		conv.SetType(object.REAL)
		x, _ := strconv.ParseFloat(val, 64)
		conv.SetData(float64(x))
	case "SHORTREAL":
		conv.SetType(object.SHORTREAL)
		x, _ := strconv.ParseFloat(val, 32)
		conv.SetData(float32(x))
	case "SET":
		conv.SetType(object.SET)
		x, _ := strconv.ParseInt(val, 2, 32)
		conv.SetData(big.NewInt(x))
	case "BOOLEAN":
		conv.SetType(object.BOOLEAN)
		if val == "TRUE" {
			conv.SetData(true)
		} else if val == "FALSE" {
			conv.SetData(false)
		} else {
			panic("wrong bool")
		}
	case "STRING":
		conv.SetType(object.STRING)
		conv.SetData(val)
	case "SHORTSTRING":
		conv.SetType(object.SHORTSTRING)
		conv.SetData(val)
	case "":
		conv.SetType(object.NOTYPE)
	default:
		panic(fmt.Sprintln("no such constant type", typ))
	}
}

var nodeMap map[string]node.Node
var objectMap map[string]object.Object

func (r *Result) doType(n *Node) (ret object.ComplexType) {
	switch n.Data.Typ.Form {
	case "BASIC":
		switch n.Data.Typ.Typ {
		case "PROCEDURE":
			t := object.NewBasicType(object.PROCEDURE)
			link := r.findLink(n, "link")
			if link != nil {
				t.SetLink(r.doObject(link))
				assert.For(t.Link() != nil, 40)
			}
			ret = t
		case "CHAR", "SHORTCHAR", "INTEGER", "LONGINT", "BYTE",
			"SHORTINT", "BOOLEAN", "REAL", "SHORTREAL", "SET":
		default:
			fmt.Println("unknown basic type", n.Data.Typ.Typ)
		}
	case "DYNAMIC":
		switch n.Data.Typ.Base {
		case "CHAR":
			n := object.NewDynArrayType(object.CHAR)
			ret = n
		case "BYTE":
			n := object.NewDynArrayType(object.BYTE)
			ret = n
		case "SHORTCHAR":
			n := object.NewDynArrayType(object.SHORTCHAR)
			ret = n
		default:
			panic(fmt.Sprintln("unknown dyn type", n.Data.Typ.Typ))
		}
	case "ARRAY":
		switch n.Data.Typ.Base {
		case "CHAR":
			n := object.NewArrayType(object.CHAR, int64(n.Data.Typ.Par))
			ret = n
		default:
			panic(fmt.Sprintln("unknown array type", n.Data.Typ.Typ))
		}
	case "RECORD":
		switch n.Data.Typ.Base {
		case "NOTYP":
			n := object.NewRecordType(n.Data.Typ.Name)
			ret = n
		default:
			n := object.NewRecordType(n.Data.Typ.Name, n.Data.Typ.Base)
			ret = n
		}
		link := r.findLink(n, "link")
		if link != nil {
			ret.SetLink(r.doObject(link))
			assert.For(ret.Link() != nil, 40)
		}
		base := r.findLink(n, "base")
		if base != nil {
			ret.(object.RecordType).SetBase(r.doType(base))
			assert.For(ret.(object.RecordType).BaseType() != nil, 41)
		}
	default:
		panic(fmt.Sprintln("unknown form", n.Data.Typ.Form))
	}
	return ret
}

func (r *Result) doObject(n *Node) object.Object {
	assert.For(n != nil, 20)
	var ret object.Object
	ret = objectMap[n.Id]
	if ret == nil {
		switch n.Data.Obj.Mode {
		case "head":
			ret = object.New(object.HEAD)
		case "variable":
			ret = object.New(object.VARIABLE)
			initType(n.Data.Obj.Typ, ret.(object.VariableObject))
		case "local procedure":
			ret = object.New(object.LOCAL_PROC)
			ret.SetType(object.PROCEDURE)
		case "external procedure":
			ret = object.New(object.EXTERNAL_PROC)
			ret.SetType(object.PROCEDURE)
		case "constant":
			ret = object.New(object.CONSTANT)
			convertData(n.Data.Obj.Typ, n.Data.Obj.Value, ret.(object.ConstantObject))
			//fmt.Println(n.Data.Obj.Name, " ", ret.(object.ConstantObject).Data())
		case "parameter":
			ret = object.New(object.PARAMETER)
			initType(n.Data.Obj.Typ, ret.(object.ParameterObject))
		case "field":
			ret = object.New(object.FIELD)
			initType(n.Data.Obj.Typ, ret.(object.FieldObject))
		case "type":
			ret = object.New(object.TYPE)
		default:
			fmt.Println(n.Data.Obj.Mode)
			panic("no such object mode")
		}
	}
	if ret != nil {
		objectMap[n.Id] = ret
		ret.SetName(n.Data.Obj.Name)

		link := r.findLink(n, "link")
		if link != nil {
			ret.SetLink(r.doObject(link))
			if ret.Link() == nil {
				panic("error in object")
			}
		}

		typ := r.findLink(n, "type")
		if typ != nil {
			ret.SetComplex(r.doType(typ))
			if ret.Complex() != nil {
				ret.SetType(object.COMPLEX)
			} else {
				//fmt.Println("not a complex type")
			}
		}

	}
	return ret
}

func (r *Result) buildScope(list []Node) (ro []object.Object, rt []object.ComplexType) {
	if list == nil {
		return nil, nil
	}
	ro = make([]object.Object, 0)
	rt = make([]object.ComplexType, 0)
	for i := range list {
		switch {
		case list[i].Data.Obj != nil:
			obj := r.doObject(&list[i])
			if obj != nil {
				ro = append(ro, obj)
			}
		case list[i].Data.Typ != nil:
			typ := r.doType(&list[i])
			if typ != nil {
				rt = append(rt, typ)
			}
		default:
			panic("no such object type")
		}

	}
	return ro, rt
}

func (r *Result) buildNode(n *Node) (ret node.Node) {
	assert.For(n != nil, 20)
	ret = nodeMap[n.Id]
	var proc node.ProcedureNode
	if ret == nil {
		switch n.Data.Nod.Class {
		case "enter":
			ret = node.New(constant.ENTER)
			switch n.Data.Nod.Enter {
			case "module":
				ret.(node.EnterNode).SetEnter(enter.MODULE)
			case "procedure":
				ret.(node.EnterNode).SetEnter(enter.PROCEDURE)
			default:
				panic("no such enter type")
			}
		case "variable":
			ret = node.New(constant.VARIABLE)
		case "dyadic":
			ret = node.New(constant.DYADIC)
			switch n.Data.Nod.Operation {
			case "+":
				ret.(node.OperationNode).SetOperation(operation.PLUS)
			case "-":
				ret.(node.OperationNode).SetOperation(operation.MINUS)
			case "=":
				ret.(node.OperationNode).SetOperation(operation.EQUAL)
			case "<":
				ret.(node.OperationNode).SetOperation(operation.LESSER)
			case "<=":
				ret.(node.OperationNode).SetOperation(operation.LESS_EQUAL)
			case "len":
				ret.(node.OperationNode).SetOperation(operation.LEN)
			case "#":
				ret.(node.OperationNode).SetOperation(operation.NOT_EQUAL)
			case ">":
				ret.(node.OperationNode).SetOperation(operation.GREATER)
			case "max":
				ret.(node.OperationNode).SetOperation(operation.MAX)
			case "min":
				ret.(node.OperationNode).SetOperation(operation.MIN)
			default:
				panic(fmt.Sprintln("no such operation", n.Data.Nod.Operation))
			}
		case "constant":
			ret = node.New(constant.CONSTANT)
			convertData(n.Data.Nod.Typ, n.Data.Nod.Value, ret.(node.ConstantNode))
			//fmt.Println(ret.(node.ConstantNode).Data())
			x := ret.(node.ConstantNode)
			if n.Data.Nod.Min != nil {
				val, _ := strconv.Atoi(n.Data.Nod.Min.X)
				x.SetMin(val)
			}
			if n.Data.Nod.Max != nil {
				val, _ := strconv.Atoi(n.Data.Nod.Max.X)
				x.SetMax(val)
			}
		case "assign":
			ret = node.New(constant.ASSIGN)
			switch n.Data.Nod.Statement {
			case "assign":
				ret.(node.AssignNode).SetStatement(statement.ASSIGN)
			case "inc":
				ret.(node.AssignNode).SetStatement(statement.INC)
			case "dec":
				ret.(node.AssignNode).SetStatement(statement.DEC)
			case "incl":
				ret.(node.AssignNode).SetStatement(statement.INCL)
			case "excl":
				ret.(node.AssignNode).SetStatement(statement.EXCL)
			default:
				panic(fmt.Sprintln("unknown assign statement", n.Data.Nod.Statement))
			}
		case "call":
			ret = node.New(constant.CALL)
		case "procedure":
			ret = node.New(constant.PROCEDURE)
			proc = ret.(node.ProcedureNode)
		case "parameter":
			ret = node.New(constant.PARAMETER)
		case "return":
			ret = node.New(constant.RETURN)
		case "monadic":
			ret = node.New(constant.MONADIC)
			switch n.Data.Nod.Operation {
			case "convert":
				ret.(node.OperationNode).SetOperation(operation.CONVERT)
				initType(n.Data.Nod.Typ, ret.(node.MonadicNode))
			case "not":
				ret.(node.OperationNode).SetOperation(operation.NOT)
			case "is":
				ret.(node.OperationNode).SetOperation(operation.IS)
			case "abs":
				ret.(node.OperationNode).SetOperation(operation.ABS)
			case "odd":
				ret.(node.OperationNode).SetOperation(operation.ODD)
			case "cap":
				ret.(node.OperationNode).SetOperation(operation.CAP)
			case "bits":
				ret.(node.OperationNode).SetOperation(operation.BITS)
			default:
				panic("no such operation " + n.Data.Nod.Operation)
			}
		case "conditional":
			ret = node.New(constant.CONDITIONAL)
		case "if":
			ret = node.New(constant.IF)
		case "while":
			ret = node.New(constant.WHILE)
		case "repeat":
			ret = node.New(constant.REPEAT)
		case "loop":
			ret = node.New(constant.LOOP)
		case "exit":
			ret = node.New(constant.EXIT)
		case "dereferencing":
			ret = node.New(constant.DEREF)
		case "field":
			ret = node.New(constant.FIELD)
		case "init":
			ret = node.New(node.INIT)
		case "index":
			ret = node.New(constant.INDEX)
		case "trap":
			ret = node.New(constant.TRAP)
		case "with":
			ret = node.New(constant.WITH)
		case "guard":
			ret = node.New(constant.GUARD)
			typ := r.findLink(n, "type")
			if typ != nil {
				ret.(node.GuardNode).SetType(r.doType(typ))
				if ret.(node.GuardNode).Type() == nil {
					panic("error in node")
				}
			}
		case "case":
			ret = node.New(constant.CASE)
		case "else":
			ret = node.New(constant.ELSE)
			x := ret.(node.ElseNode)
			if n.Data.Nod.Min != nil {
				val, _ := strconv.Atoi(n.Data.Nod.Min.X)
				x.Min(val)
			}
			if n.Data.Nod.Max != nil {
				val, _ := strconv.Atoi(n.Data.Nod.Max.X)
				x.Max(val)
			}
		case "do":
			ret = node.New(constant.DO)
		case "range":
			ret = node.New(constant.RANGE)
		default:
			fmt.Println(n.Data.Nod.Class)
			panic("no such node type")
		}
	}
	if ret != nil {
		nodeMap[n.Id] = ret
		left := r.findLink(n, "left")
		if left != nil {
			ret.SetLeft(r.buildNode(left))
			if ret.Left() == nil {
				panic("error in node")
			}
		}
		right := r.findLink(n, "right")
		if right != nil {
			ret.SetRight(r.buildNode(right))
			if ret.Right() == nil {
				panic("error in node")
			}
		}
		link := r.findLink(n, "link")
		if link != nil {
			ret.SetLink(r.buildNode(link))
			if ret.Link() == nil {
				panic("error in node")
			}
		}
		object := r.findLink(n, "object")
		if object != nil {
			ret.SetObject(r.doObject(object))
			if ret.Object() == nil {
				panic("error in node")
			}
		} else {
			assert.For(proc == nil, 60) //у процедуры просто не может не быть объекта
		}

	}
	return ret
}

func buildMod(r *Result) (nodeList []node.Node, scopeList map[node.Node][]object.Object, typeList map[node.Node][]object.ComplexType, root node.Node) {
	scopes := make(map[string][]object.Object, 0)
	types := make(map[string][]object.ComplexType, 0)
	for i := range r.GraphList {
		if r.GraphList[i].CptScope != "" {
			scopes[r.GraphList[i].CptScope], types[r.GraphList[i].CptScope] = r.buildScope(r.GraphList[i].NodeList)
		}
	}
	scopeList = make(map[node.Node][]object.Object, 0)
	typeList = make(map[node.Node][]object.ComplexType, 0)
	for i := range r.GraphList {
		if r.GraphList[i].CptProc != "" {
			nodeList = make([]node.Node, 0)
			for j := range r.GraphList[i].NodeList {
				node := &r.GraphList[i].NodeList[j]
				ret := r.buildNode(node)
				nodeList = append(nodeList, ret)
				if scopes[node.Id] != nil {
					scopeList[ret] = scopes[node.Id]
				}
				if types[node.Id] != nil {
					typeList[ret] = types[node.Id]
				}
				if (node.Data.Nod.Class == "enter") && (node.Data.Nod.Enter == "module") {
					root = ret
				}
			}
		}
	}
	return nodeList, scopeList, typeList, root
}

func DoAST(r *Result) (mod *module.Module) {
	nodeMap = make(map[string]node.Node)
	objectMap = make(map[string]object.Object)
	mod = new(module.Module)
	mod.Nodes, mod.Objects, mod.Types, mod.Enter = buildMod(r)
	fmt.Println(len(mod.Nodes), len(mod.Objects))
	nodeMap = nil
	objectMap = nil
	return mod
}
