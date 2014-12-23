package xev

import (
	"fmt"
	"fw/cp/constant"
	"fw/cp/constant/enter"
	"fw/cp/constant/operation"
	"fw/cp/module"
	"fw/cp/node"
	"fw/cp/object"
	"fw/cp/statement"
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
		x, _ := strconv.ParseInt(val, 16, 32)
		conv.SetData(int32(x))
	case "SHORTINT":
		conv.SetType(object.SHORTINT)
		x, _ := strconv.ParseInt(val, 16, 16)
		conv.SetData(int16(x))
	case "LONGINT":
		conv.SetType(object.LONGINT)
		x, _ := strconv.ParseInt(val, 16, 64)
		conv.SetData(x)
	case "BYTE":
		conv.SetType(object.BYTE)
		x, _ := strconv.ParseInt(val, 16, 8)
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
	case "":
		conv.SetType(object.NOTYPE)
	default:
		panic("no such constant type")
	}
}

var nodeMap map[string]node.Node
var objectMap map[string]object.Object

func (r *Result) buildObject(n *Node) object.Object {
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
			ret = object.New(object.LOCAL_PROCEDURE)
		case "constant":
			ret = object.New(object.CONSTANT)
			convertData(n.Data.Obj.Typ, n.Data.Obj.Value, ret.(object.ConstantObject))
			//fmt.Println(n.Data.Obj.Name, " ", ret.(object.ConstantObject).Data())
		case "parameter":
			ret = object.New(object.PARAMETER)
			initType(n.Data.Obj.Typ, ret.(object.ParameterObject))
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
			ret.SetLink(r.buildObject(link))
			if ret.Link() == nil {
				panic("error in object")
			}
		}

	}
	return ret
}

func (r *Result) buildObjectList(list []Node) []object.Object {
	assert.For(list != nil, 20)
	ret := make([]object.Object, 0)
	for i := range list {
		obj := r.buildObject(&list[i])
		if obj != nil {
			ret = append(ret, obj)
		}
	}

	return ret
}

func (r *Result) buildNode(n *Node) (ret node.Node) {
	assert.For(n != nil, 20)
	ret = nodeMap[n.Id]
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
			case "plus":
				ret.(node.OperationNode).SetOperation(operation.PLUS)
			case "minus":
				ret.(node.OperationNode).SetOperation(operation.MINUS)
			case "equal":
				ret.(node.OperationNode).SetOperation(operation.EQUAL)
			case "lesser":
				ret.(node.OperationNode).SetOperation(operation.LESSER)
			default:
				panic("no such operation")
			}
		case "constant":
			ret = node.New(constant.CONSTANT)
			convertData(n.Data.Nod.Typ, n.Data.Nod.Value, ret.(node.ConstantNode))
			//fmt.Println(ret.(node.ConstantNode).Data())
		case "assign":
			ret = node.New(constant.ASSIGN)
			switch n.Data.Nod.Statement {
			case "assign":
				ret.(node.AssignNode).SetStatement(statement.ASSIGN)
			default:
				panic("unknown assign statement")
			}
		case "call":
			ret = node.New(constant.CALL)
		case "procedure":
			ret = node.New(constant.PROCEDURE)
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
			default:
				panic("no such operation")
			}
		case "conditional":
			ret = node.New(constant.CONDITIONAL)
		case "if":
			ret = node.New(constant.IF)
		case "while":
			ret = node.New(constant.WHILE)
		case "repeat":
			ret = node.New(constant.REPEAT)
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
			ret.SetObject(r.buildObject(object))
			if ret.Object() == nil {
				panic("error in node")
			}
		}
	}
	return ret
}

func buildMod(r *Result) (nodeList []node.Node, scopeList map[node.Node][]object.Object, root node.Node) {
	scopes := make(map[string][]object.Object, 0)
	for i := range r.GraphList {
		if r.GraphList[i].CptScope != "" {
			scopes[r.GraphList[i].CptScope] = r.buildObjectList(r.GraphList[i].NodeList)
		}
	}
	scopeList = make(map[node.Node][]object.Object, 0)
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
				if (node.Data.Nod.Class == "enter") && (node.Data.Nod.Enter == "module") {
					root = ret
				}
			}
		}
	}
	return nodeList, scopeList, root
}

func DoAST(r *Result) (mod *module.Module) {
	nodeMap = make(map[string]node.Node)
	objectMap = make(map[string]object.Object)
	mod = new(module.Module)
	mod.Nodes, mod.Objects, mod.Enter = buildMod(r)
	fmt.Println(len(mod.Nodes), len(mod.Objects))
	nodeMap = nil
	objectMap = nil
	return mod
}
