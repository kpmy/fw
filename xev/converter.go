package xev

import (
	"cp/constant"
	"cp/constant/enter"
	"cp/constant/operation"
	"cp/node"
	"cp/object"
	"cp/statement"
	"fmt"
	"strconv"
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

var nodeMap map[string]node.Node
var objectMap map[string]object.Object

func (r *Result) buildObject(n *Node) object.Object {
	if n == nil {
		panic("n is nil")
	}
	var ret object.Object
	ret = objectMap[n.Id]
	if ret == nil {
		switch n.Data.Obj.Mode {
		case "head":
			ret = object.New(object.HEAD)
		case "variable":
			ret = object.New(object.VARIABLE)
		case "local procedure":
			ret = object.New(object.LOCAL_PROCEDURE)
		default:
			panic("no such object mode")
		}
	}
	if ret != nil {
		objectMap[n.Id] = ret
		ret.SetName(n.Data.Obj.Name)
		switch n.Data.Obj.Typ {
		case "":
			ret.SetType(object.NOTYPE)
		case "INTEGER":
			ret.SetType(object.INTEGER)
		default:
			panic("no such object type")
		}
	}
	return ret
}

func (r *Result) buildObjectList(list []Node) []object.Object {
	if list == nil {
		panic("list is nil")
	}
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
	if n == nil {
		panic("n is nil")
	}
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
			default:
				panic("no such operation")
			}
		case "constant":
			ret = node.New(constant.CONSTANT)
			switch n.Data.Nod.Typ {
			case "INTEGER":
				ret.(node.ConstantNode).SetType(object.INTEGER)
				x, _ := strconv.Atoi(n.Data.Nod.Value)
				ret.(node.ConstantNode).SetData(x)
			default:
				panic("no such constant type")
			}
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
		default:
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

func buildMod(r *Result) (nodeList []node.Node, objList []object.Object, root node.Node) {
	for i := range r.GraphList {
		if r.GraphList[i].CptScope != "" {
			objList = r.buildObjectList(r.GraphList[i].NodeList)
		}
		if r.GraphList[i].CptProc != "" {
			nodeList = make([]node.Node, 0)
			for j := range r.GraphList[i].NodeList {
				node := &r.GraphList[i].NodeList[j]
				ret := r.buildNode(&r.GraphList[i].NodeList[j])
				nodeList = append(nodeList, ret)
				if (node.Data.Nod.Class == "enter") && (node.Data.Nod.Enter == "module") {
					root = ret
				}
			}
		}
	}
	return nodeList, objList, root
}

func DoAST(r *Result) (ent node.Node, obj []object.Object) {
	nodeMap = make(map[string]node.Node)
	objectMap = make(map[string]object.Object)
	_, obj, ent = buildMod(r)
	fmt.Println(len(nodeMap), len(objectMap))
	nodeMap = nil
	objectMap = nil
	return ent, obj
}
