package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatementAst(t *testing.T) {
	n1 := newTokenNode(NODE_VARIABLE, "node1", nil)
	n2 := newTokenNode(NODE_VERB, "use", nil)
	root := newTokenNode(NODE_STATEMENT, "", n1, n2)
	assert.Equal(t, n1, root.Left())
	assert.Equal(t, n2, root.Right())
}
