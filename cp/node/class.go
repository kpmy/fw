package node

import (
	"fmt"
	"fw/cp"
	"fw/cp/constant"
	"fw/cp/constant/enter"
	"fw/cp/constant/operation"
	"fw/cp/constant/statement"
	"fw/cp/object"
	"ypk/assert"
)

const INIT constant.Class = -1
const COMPOUND constant.Class = -2

type InitNode interface {
	Node
	Statement
	self() InitNode
}

type CompNode interface {
	Node
	Statement
	self() CompNode
}

func New(class constant.Class, id int) (ret Node) {
	switch class {
	case constant.ENTER:
		ret = new(enterNode)
	case constant.ASSIGN:
		ret = new(assignNode)
	case constant.VARIABLE:
		ret = new(variableNode)
	case constant.DYADIC:
		ret = new(dyadicNode)
	case constant.CONSTANT:
		ret = new(constantNode)
	case constant.CALL:
		ret = new(callNode)
	case constant.PROCEDURE:
		ret = new(procedureNode)
	case constant.PARAMETER:
		ret = new(parameterNode)
	case constant.RETURN:
		ret = new(returnNode)
	case constant.MONADIC:
		ret = new(monadicNode)
	case constant.CONDITIONAL:
		ret = new(conditionalNode)
	case constant.IF:
		ret = new(ifNode)
	case constant.REPEAT:
		ret = new(repeatNode)
	case constant.WHILE:
		ret = new(whileNode)
	case constant.EXIT:
		ret = new(exitNode)
	case constant.LOOP:
		ret = new(loopNode)
	case constant.DEREF:
		ret = new(derefNode)
	case constant.FIELD:
		ret = new(fieldNode)
	case INIT:
		ret = new(initNode)
	case constant.INDEX:
		ret = new(indexNode)
	case constant.TRAP:
		ret = new(trapNode)
	case constant.WITH:
		ret = new(withNode)
	case constant.GUARD:
		ret = new(guardNode)
	case constant.CASE:
		ret = new(caseNode)
	case constant.ELSE:
		ret = new(elseNode)
	case constant.DO:
		ret = new(doNode)
	case constant.RANGE:
		ret = new(rangeNode)
	case COMPOUND:
		ret = new(compNode)
	default:
		panic("no such class")
	}
	ret.Adr(cp.Next(id))
	return ret
}

type nodeFields struct {
	left, right, link Node
	obj               object.Object
	adr               cp.ID
}

type stmt struct {
	Statement
}

func (s *stmt) s() Statement { return s }

type expr struct {
	Expression
}

func (e *expr) e() Expression { return e }

type design struct {
	Designator
}

func (x *design) d() Designator { return x }

func (nf *nodeFields) Adr(a ...cp.ID) cp.ID {
	assert.For(len(a) <= 1, 20)
	if len(a) == 1 {
		nf.adr = a[0]
	}
	return nf.adr
}

func (nf *nodeFields) SetLeft(n Node) { nf.left = n }

func (nf *nodeFields) SetRight(n Node) { nf.right = n }

func (nf *nodeFields) SetLink(n Node) { nf.link = n }

func (nf *nodeFields) SetObject(o object.Object) { nf.obj = o; o.SetRef(nf) }

func (nf *nodeFields) Left() Node { return nf.left }

func (nf *nodeFields) Right() Node { return nf.right }

func (nf *nodeFields) Link() Node { return nf.link }

func (nf *nodeFields) Object() object.Object { return nf.obj }

type enterNode struct {
	nodeFields
	stmt
	enter enter.Enter
}

func (e *enterNode) SetEnter(enter enter.Enter) { e.enter = enter }
func (e *enterNode) Enter() enter.Enter         { return e.enter }

type constantNode struct {
	nodeFields
	expr
	typ      object.Type
	data     interface{}
	min, max *int
}

func (c *constantNode) SetType(t object.Type) { c.typ = t }

func (c *constantNode) SetData(data interface{}) { c.data = data }

func (c *constantNode) Data() interface{} { return c.data }

func (c *constantNode) Type() object.Type { return c.typ }

func (c *constantNode) SetMax(x int) { c.max = new(int); *c.max = x }
func (c *constantNode) Max() *int    { return c.max }
func (c *constantNode) SetMin(x int) { c.min = new(int); *c.min = x }
func (c *constantNode) Min() *int    { return c.min }

type dyadicNode struct {
	nodeFields
	operation operation.Operation
	expr
}

func (d *dyadicNode) SetOperation(op operation.Operation) { d.operation = op }

func (d *dyadicNode) Operation() operation.Operation { return d.operation }

func (d *dyadicNode) self() DyadicNode { return d }

func (d *dyadicNode) String() string {
	return fmt.Sprintln("dyadic ", d.Adr(), d.operation)
}

type assignNode struct {
	nodeFields
	stat statement.Statement
	stmt
}

func (a *assignNode) self() AssignNode { return a }

func (a *assignNode) SetStatement(s statement.Statement) { a.stat = s }

func (a *assignNode) Statement() statement.Statement { return a.stat }

type variableNode struct {
	nodeFields
	design
	expr
}

func (v *variableNode) self() VariableNode { return v }

type callNode struct {
	nodeFields
	stmt
	expr
}

func (v *callNode) self() CallNode { return v }

func (c *callNode) String() string {
	return fmt.Sprintln(c.Adr(), c.Left().Object().Name())
}

type procedureNode struct {
	nodeFields
	super bool
	design
}

func (v *procedureNode) self() ProcedureNode { return v }
func (v *procedureNode) Super(x ...string) bool {
	if len(x) > 0 {
		v.super = x[0] == "super"
	}
	return v.super
}

type parameterNode struct {
	nodeFields
	design
}

func (v *parameterNode) self() ParameterNode { return v }

type returnNode struct {
	nodeFields
	stmt
}

func (v *returnNode) self() ReturnNode { return v }

type monadicNode struct {
	nodeFields
	operation operation.Operation
	typ       object.Type
	comp      object.ComplexType
	expr
}

func (v *monadicNode) self() MonadicNode { return v }

func (v *monadicNode) SetOperation(op operation.Operation) { v.operation = op }

func (v *monadicNode) Operation() operation.Operation { return v.operation }

func (v *monadicNode) SetType(t object.Type) { v.typ = t }
func (v *monadicNode) Type() object.Type     { return v.typ }

func (v *monadicNode) Complex(x ...object.ComplexType) object.ComplexType {
	if len(x) > 0 {
		v.comp = x[0]
	}
	return v.comp
}

func (v *monadicNode) String() string {
	return fmt.Sprintln("monadic ", v.Adr(), v.operation)
}

type conditionalNode struct {
	nodeFields
	stmt
}

func (v *conditionalNode) self() ConditionalNode { return v }

type ifNode struct {
	nodeFields
}

func (v *ifNode) self() IfNode { return v }

type whileNode struct {
	nodeFields
	stmt
}

func (v *whileNode) self() WhileNode { return v }

type repeatNode struct {
	nodeFields
	stmt
}

func (v *repeatNode) self() RepeatNode { return v }

type exitNode struct {
	nodeFields
	stmt
}

func (v *exitNode) self() ExitNode { return v }

type loopNode struct {
	nodeFields
	stmt
}

func (v *loopNode) self() LoopNode { return v }

type derefNode struct {
	nodeFields
	ptr bool
	design
}

func (v *derefNode) self() DerefNode { return v }

func (c *derefNode) Ptr(x ...string) bool {
	if len(x) > 0 {
		c.ptr = x[0] == "ptr"
	}
	return c.ptr
}

type fieldNode struct {
	nodeFields
	design
}

func (v *fieldNode) self() FieldNode { return v }

type initNode struct {
	nodeFields
	stmt
}

func (v *initNode) self() InitNode { return v }

type compNode struct {
	nodeFields
	stmt
}

func (v *compNode) self() CompNode { return v }

type indexNode struct {
	nodeFields
	design
}

func (v *indexNode) self() IndexNode { return v }

type trapNode struct {
	nodeFields
	stmt
}

func (v *trapNode) self() TrapNode { return v }

type withNode struct {
	nodeFields
	stmt
}

func (v *withNode) self() WithNode { return v }

type guardNode struct {
	nodeFields
	typ object.ComplexType
	design
}

func (v *guardNode) self() GuardNode              { return v }
func (v *guardNode) SetType(t object.ComplexType) { v.typ = t }
func (v *guardNode) Type() object.ComplexType     { return v.typ }

type caseNode struct {
	nodeFields
	stmt
}

func (v *caseNode) self() CaseNode { return v }

type elseNode struct {
	nodeFields
	min, max int
}

func (v *elseNode) Min(x ...int) int {
	if len(x) > 0 {
		v.min = x[0]
	}
	return v.min
}

//func (c *elseNode) SetMax(x int) { c.max = x }
func (c *elseNode) Max(x ...int) int {
	if len(x) > 0 {
		c.max = x[0]
	}
	return c.max
}

//func (c *elseNode) SetMin(x int) { c.min = x }

type doNode struct {
	nodeFields
}

func (v *doNode) self() DoNode { return v }

type rangeNode struct {
	nodeFields
}

func (v *rangeNode) self() RangeNode { return v }
