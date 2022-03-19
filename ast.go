package main

// Abstract Syntax Tree

var count int

// Node is interface with methods
// that every AST node must implement
type Node interface {
	NodeId() int // identifier
	Left() Node
	Right() Node
	AppendChild(Node)
}

// BaseNode is base node with reference to childs
type BaseNode struct {
	Id        int
	currChild int
	children  []Node
}

// NodeType supported AST nodes
// These nodes are not the same as lex nodes
type NodeType int

const (
	NODE_ILLEGAL   NodeType = iota
	NODE_SUBJECT            // Subject
	NODE_STATEMENT          // Statement
	NODE_VARIABLE           // Variable
	NODE_OP                 // Operant
	NODE_VALUE              // Value
	NODE_VERB               // Verb
	NODE_RESOURCE           // Resource
	NODE_LOCATION           // Location
	NODE_CONDITION          // Condition
)

type TokenNode struct {
	BaseNode
	Token string
	NodeType
}

// CreateAstNode creates Base node
func CreateAstNode(id int) BaseNode {
	b := BaseNode{Id: count}
	b.children = make([]Node, 0)
	return b
}

// AppendChild adds new child node
func (b *BaseNode) AppendChild(n Node) {
	b.children = append(b.children, n)
}

// NodeId returns node unique idenifier
func (b *BaseNode) NodeId() int {
	return b.Id
}

// Left returns left node
func (b *BaseNode) Left() Node {
	b.currChild = 0
	if len(b.children) < 1 {
		return nil
	}
	return b.children[b.currChild]
}

// Right returns right nodes
func (b *BaseNode) Right() Node {
	b.currChild = b.currChild + 1
	if b.currChild >= len(b.children) {
		return nil
	}
	return b.children[b.currChild]
}

// newTokenNode creates new AST node with nType
func newTokenNode(nType NodeType, str string, childs ...Node) Node {
	b := CreateAstNode(count)
	for _, ch := range childs {
		b.AppendChild(ch)
	}

	e := &TokenNode{
		BaseNode: b,
		NodeType: nType,
	}
	if str != "" {
		e.Token = str
	}
	count++
	return e
}
