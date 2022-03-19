package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const dotFileName = "output.dot"

func intToString(value int) string {
	return strconv.Itoa(value)
}

// createLabel returns AST node name
func createLabel(node Node) (string, error) {
	var n, ok = node.(*TokenNode)
	if ok {
		var label string
		switch n.NodeType {
		case NODE_RESOURCE:
			label = fmt.Sprintf("Resource\n%s", n.Token)
		case NODE_VERB:
			label = fmt.Sprintf("Verb\n%s", n.Token)
		case NODE_STATEMENT:
			label = "Statement"
		case NODE_LOCATION:
			label = "Location"
		case NODE_VARIABLE:
			label = n.Token
		case NODE_OP:
			label = n.Token
		case NODE_VALUE:
			label = n.Token
		case NODE_CONDITION:
			label = "Condition"
		case NODE_SUBJECT:
			label = "Subject"
		default:
			label = "Unknown"
		}
		return label, nil
	} else {
		return "", fmt.Errorf("unexpected node type %T", node)
	}
}

// Scan scans all nodes in the tree recursively
func scan(node Node, edges *[]string, labels *[]string) error {
	if node == nil {
		return nil
	}

	for e := node.Left(); e != nil; e = node.Right() {
		edge1 := intToString(node.NodeId())
		edge2 := intToString(e.NodeId())

		label1, err := createLabel(node)
		if err != nil {
			return err
		}

		label2, err := createLabel(e)
		if err != nil {
			return err
		}

		edge := fmt.Sprintf("\t %s -> %s", edge1, edge2)
		label := fmt.Sprintf("\t %s [label=\"%s\"];\n", edge1, label1)
		label += fmt.Sprintf("\t %s [label=\"%s\"];", edge2, label2)

		*edges = append(*edges, edge)
		*labels = append(*labels, label)
	}

	for e := node.Left(); e != nil; e = node.Right() {
		scan(e, edges, labels)
	}
	return nil
}

// Convert creates DOT (graph description language) file from Node
func GenerateDotFormat(node Node, outputfile string) error {
	file, err := os.Create(outputfile)
	defer file.Close()
	if err != nil {
		return err
	}
	result := "digraph G {" + "\n"
	e := []string{}
	l := []string{}
	scan(node, &e, &l)
	result += fmt.Sprintf("%s \n %s \n}", strings.Join(e, "\n"), strings.Join(l, "\n"))
	_, err = file.WriteString(result)
	return err
}
