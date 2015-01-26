package xev

import (
	"encoding/xml"
	"fmt"
	"io"
)

type CptType struct {
	Form string `xml:"urn:bbcb:component:dev:cpt form,attr"`
	Typ  string `xml:"urn:bbcb:component:dev:cpt type,attr"`
	Base string `xml:"urn:bbcb:component:dev:cpt base,attr"`
	Par  int    `xml:"urn:bbcb:component:dev:cpt par,attr"`
	Name string `xml:"urn:bbcb:component:dev:cpt name,attr"`
}

type CptObject struct {
	Name string `xml:"urn:bbcb:component:dev:cpt name,attr"`
	Mode string `xml:"urn:bbcb:component:dev:cpt mode,attr"`
	Typ  string `xml:"urn:bbcb:component:dev:cpt type,attr"`
	//	Leaf  bool   `xml:"urn:bbcb:component:dev:cpt leaf,attr"`
	Value string `xml:",chardata"`
}

type Value struct {
	X string `xml:",chardata"`
}

type CptNode struct {
	Class string `xml:"urn:bbcb:component:dev:cpt class,attr"`
	//опциональные параметры
	Typ       string `xml:"urn:bbcb:component:dev:cpt type,attr"`
	Enter     string `xml:"urn:bbcb:component:dev:cpt enter,attr"`
	Operation string `xml:"urn:bbcb:component:dev:cpt operation,attr"`
	Value     string `xml:",chardata"`
	Min       *Value `xml:"urn:bbcb:component:dev:cpt min"`
	Max       *Value `xml:"urn:bbcb:component:dev:cpt max"`
	Statement string `xml:"urn:bbcb:component:dev:cpt statement,attr"`
	Proto     string `xml:"urn:bbcb:component:dev:cpt proto,attr"`
	From      string `xml:"urn:bbcb:component:dev:cpt from,attr"`
	Proc      string `xml:"urn:bbcb:component:dev:cpt proc,attr"`
}

type NodeData struct {
	Obj *CptObject `xml:"urn:bbcb:component:dev:cpt object"`
	Nod *CptNode   `xml:"urn:bbcb:component:dev:cpt node"`
	Typ *CptType   `xml:"urn:bbcb:component:dev:cpt type"`
}

type Node struct {
	Id   int       `xml:"id,attr"`
	Data *NodeData `xml:"data"`
}

type Edge struct {
	Target  int    `xml:"target,attr"`
	Source  int    `xml:"source,attr"`
	CptLink string `xml:"urn:bbcb:component:dev:cpt link,attr"`
}

type GraphData struct {
	CptScope string `xml:"urn:bbcb:component:dev:cpt scope,attr"`
	CptProc  string `xml:"urn:bbcb:component:dev:cpt proc,attr"`
}

type Graph struct {
	NodeList []Node `xml:"node"`
	EdgeList []Edge `xml:"edge"`
	GraphData
}

type Result struct {
	GraphList []Graph `xml:"graph"`
}

func traverseNode(n *Node) {
	fmt.Println(n.Id, n.Data)
	if n.Data.Nod != nil {
		fmt.Println(n.Data.Nod.Min)
		fmt.Println(n.Data.Nod.Max)
	}
}

func traverseEdge(e *Edge) {
	fmt.Println(e.Source, e.Target, e.CptLink)
}

func traverseGraph(g *Graph) {
	fmt.Println("scope", g.CptScope)
	fmt.Println("proc", g.CptProc)
	for n := range g.NodeList {
		fmt.Println("node", n)
		traverseNode(&g.NodeList[n])
	}
	for e := range g.EdgeList {
		fmt.Println("edge", e)
		traverseEdge(&g.EdgeList[e])
	}
}

func traverse(r *Result) {
	for g := range r.GraphList {
		fmt.Println("graph", g)
		traverseGraph(&r.GraphList[g])
	}
}

func LoadOXF(rd io.Reader) *Result {
	r := new(Result)
	dec := xml.NewDecoder(rd)
	err := dec.Decode(r)
	if err == nil {
		//fmt.Println(len(r.GraphList))
		//traverse(r)
	} else {
		panic("xml parse error")
	}
	return r
}
