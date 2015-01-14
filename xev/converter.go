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

var ncache map[int]*Node

func (r *Result) findNode(id int) (ret *Node) {
	//fmt.Print("|")
	ret = ncache[id]
	if ret == nil {
		for i := 0; i < len(r.GraphList) && (ret == nil); i++ {
			for j := 0; j < len(r.GraphList[i].NodeList) && (ret == nil); j++ {
				if id == r.GraphList[i].NodeList[j].Id {
					ret = &r.GraphList[i].NodeList[j]
				}
			}
		}
		if ret != nil {
			ncache[id] = ret
		}
	}
	return ret
}

type eid struct {
	id   int
	link string
}

var ecache map[eid]*Node

func (r *Result) findLink(n *Node, ls string) (ret *Node) {
	//fmt.Print("-")
	ret = ecache[eid{id: n.Id, link: ls}]
	if ret == nil {
		target := -1
		for i := range r.GraphList {
			for j := range r.GraphList[i].EdgeList {
				if (r.GraphList[i].EdgeList[j].Source == n.Id) && (r.GraphList[i].EdgeList[j].CptLink == ls) {
					target = r.GraphList[i].EdgeList[j].Target
				}
			}
		}
		if target != -1 {
			ret = r.findNode(target)
		}
		if ret != nil {
			ecache[eid{id: n.Id, link: ls}] = ret
		}
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
		panic(fmt.Sprintln("no such type", typ))
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

var nodeMap map[int]node.Node
var objectMap map[int]object.Object
var typeMap map[int]object.ComplexType

func reset() {
	nodeMap = make(map[int]node.Node)
	objectMap = make(map[int]object.Object)
	typeMap = make(map[int]object.ComplexType)
	ncache = make(map[int]*Node)
	ecache = make(map[eid]*Node)
}

func init() { reset() }

func (r *Result) doType(n *Node) (ret object.ComplexType) {
	//fmt.Println("type", n.Id)
	ret = typeMap[n.Id]
	if ret == nil {
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
			case "POINTER":
				t := object.NewPointerType(n.Data.Typ.Name)
				base := r.findLink(n, "base")
				if base != nil {
					t.Base(r.doType(base))
					assert.For(t.Base() != nil, 41)
				}
				ret = t
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
	}
	if ret != nil {
		typeMap[n.Id] = ret
	}
	return ret
}

func (r *Result) doObject(n *Node) (ret object.Object) {
	//fmt.Println("object", n.Id)
	assert.For(n != nil, 20)
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
		case "type procedure":
			ret = object.New(object.TYPE_PROC)
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
		case "module":
			ret = object.New(object.MODULE)
		default:
			fmt.Println(n.Data.Obj.Mode)
			panic("no such object mode")
		}
	}
	if ret != nil && objectMap[n.Id] == nil {
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
	//fmt.Println("node", n.Id)
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
			ret.(node.OperationNode).SetOperation(operation.This(n.Data.Nod.Operation))
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
			ret.(node.AssignNode).SetStatement(statement.This(n.Data.Nod.Statement))
		case "call":
			ret = node.New(constant.CALL)
		case "procedure":
			ret = node.New(constant.PROCEDURE)
			proc = ret.(node.ProcedureNode)
			proc.Super(n.Data.Nod.Proc)
		case "parameter":
			ret = node.New(constant.PARAMETER)
		case "return":
			ret = node.New(constant.RETURN)
		case "monadic":
			ret = node.New(constant.MONADIC)
			ret.(node.OperationNode).SetOperation(operation.This(n.Data.Nod.Operation))
			switch n.Data.Nod.Operation {
			case "CONV":
				ret.(node.OperationNode).SetOperation(operation.ALIEN_CONV)
				initType(n.Data.Nod.Typ, ret.(node.MonadicNode))
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
			ret.(node.DerefNode).Ptr(n.Data.Nod.From)
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
		case "compound":
			ret = node.New(node.COMPOUND)
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
			//pk, 20150112, у процедуры из другого модуля - может и не быть объекта
			//assert.For(proc == nil, 60) //у процедуры просто не может не быть объекта
		}

	}
	return ret
}

func buildMod(r *Result) *module.Module {
	//временные структуры создаем по очереди, чтобы корректно заполнять все ссылки на объекты/узлы
	type scope struct {
		mod    string
		scopes map[int][]object.Object
		types  map[int][]object.ComplexType
	}
	scopes := make(map[int]*scope, 0)
	for _, g := range r.GraphList {
		if g.CptScope != "" {
			sc, _ := strconv.Atoi(g.CptScope)
			var imp int
			if sc >= 0 {
				imp = 0
			} else {
				imp = sc
			}
			this := scopes[imp]
			if this == nil {
				this = &scope{}
				this.scopes = make(map[int][]object.Object, 0)
				this.types = make(map[int][]object.ComplexType, 0)
				scopes[imp] = this
			}
			if this.mod == "" {
				this.mod = g.CptProc
			}
			this.scopes[sc], this.types[sc] = r.buildScope(g.NodeList)
			fmt.Println(sc, len(this.scopes[sc]), len(this.types[sc]))
		}
	}
	//временные структуры перегоняем в рабочие
	var (
		nodeList  []node.Node                        = make([]node.Node, 0)
		scopeList map[node.Node][]object.Object      = make(map[node.Node][]object.Object, 0)
		typeList  map[node.Node][]object.ComplexType = make(map[node.Node][]object.ComplexType, 0)
		impList   map[string]module.Import           = make(map[string]module.Import, 0)
		root      node.Node
	)
	for _, g := range r.GraphList {
		if g.CptScope == "" {
			for _, nl := range g.NodeList {
				node := &nl
				ret := r.buildNode(node)
				nodeList = append(nodeList, ret)
				if scopes[0].scopes[node.Id] != nil {
					scopeList[ret] = scopes[0].scopes[node.Id]
				}
				if scopes[0].types[node.Id] != nil {
					typeList[ret] = scopes[0].types[node.Id]
				}
				if (node.Data.Nod.Class == "enter") && (node.Data.Nod.Enter == "module") {
					root = ret
				}
			}
		}
	}
	for k, v := range scopes {
		if k < 0 {
			impList[v.mod] = module.Import{Objects: v.scopes[k], Name: v.mod}
			fmt.Println("типы не учтены", len(v.types[k]))
		}
	}
	return &module.Module{Nodes: nodeList, Objects: scopeList, Types: typeList, Enter: root, Imports: impList}
}

func DoAST(r *Result) (mod *module.Module) {
	mod = buildMod(r)
	fmt.Println(len(mod.Nodes), len(mod.Objects))
	reset()
	return mod
}
