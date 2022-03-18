package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const chunkSize = 64000

func TestPlotDot(t *testing.T) {
	testFileName := "output-test.dot"
	goldenFileName := "testdata/golden.file"
	n1 := newTokenNode(NODE_VARIABLE, "variable1", nil)
	n2 := newTokenNode(NODE_VERB, "read", nil)
	root := newStatementNode(n1, n2, nil, nil, nil)
	err := GenerateDotFormat(root, testFileName)
	if err != nil {
		t.Fatal(err)
	}
	// compare files
	assert.True(t, deepCompare(goldenFileName, testFileName))
}

// Used from https://stackoverflow.com/questions/29505089/how-can-i-compare-two-files-in-golang
func deepCompare(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
