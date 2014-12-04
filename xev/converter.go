package xev

import (
	"cp/node"
	"cp/object"
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

func (r *Result) buildNode(n *Node) node.Node {
	if n == nil {
		panic("n is nil")
	}
	var ret node.Node
	ret = nodeMap[n.Id]
	if ret == nil {
		switch n.Data.Nod.Class {
		case "enter":
			ret = node.New(node.ENTER)
			switch n.Data.Nod.Enter {
			case "module":
				ret.(node.EnterNode).SetEnter(node.MODULE)
			default:
				panic("no such enter type")
			}
		case "variable":
			ret = node.New(node.VARIABLE)
		case "dyadic":
			ret = node.New(node.DYADIC)
			switch n.Data.Nod.Operation {
			case "plus":
				ret.(node.OperationNode).SetOperation(node.PLUS)
			default:
				panic("no such operation")
			}
		case "constant":
			ret = node.New(node.CONSTANT)
			switch n.Data.Nod.Typ {
			case "INTEGER":
				ret.(node.ConstantNode).SetType(object.INTEGER)
				x, _ := strconv.Atoi(n.Data.Nod.Value)
				ret.(node.ConstantNode).SetData(x)
			default:
				panic("no such constant type")
			}
		case "assign":
			ret = node.New(node.ASSIGN)
		default:
			panic("no such node type")
		}
	}
	if ret != nil {
		nodeMap[n.Id] = ret
		left := r.findLink(n, "left")
		if left != nil {
			ret.SetLeft(r.buildNode(left))
		}
		right := r.findLink(n, "right")
		if right != nil {
			ret.SetRight(r.buildNode(right))
		}
		link := r.findLink(n, "link")
		if link != nil {
			ret.SetLink(r.buildNode(link))
		}
		object := r.findLink(n, "object")
		if object != nil {
			ret.SetObject(r.buildObject(object))
		}
	}
	return ret
}

func buildMod(r *Result) (node.Node, []object.Object) {
	var root node.Node
	var list []object.Object
	for i := range r.GraphList {
		if r.GraphList[i].CptScope != "" {
			list = r.buildObjectList(r.GraphList[i].NodeList)
		}
		if r.GraphList[i].CptProc != "" {
			for j := range r.GraphList[i].NodeList {
				root = r.buildNode(&r.GraphList[i].NodeList[j])
			}
		}
	}
	return root, list
}

func DoAST(r *Result) {
	nodeMap = make(map[string]node.Node)
	objectMap = make(map[string]object.Object)
	_, _ = buildMod(r)
	fmt.Println(len(nodeMap), len(objectMap))
	nodeMap = nil
	objectMap = nil
}
