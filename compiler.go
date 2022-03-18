package main

import (
	"fmt"
	"io"
)

type Compiler struct {
	err error
	ast Node
}

type SyntaxError struct {
	Row int
	Col int
	Msg string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Row, e.Col, e.Msg)
}

// NewCompiler initializes compiler
func NewCompiler() *Compiler {
	return &Compiler{}
}

// SetAstRoot sets root to Abstract Source Tree
func (c *Compiler) SetAstRoot(root Node) {
	c.ast = root
}

// Parse parses statement from io.Reader
func (c *Compiler) Parse(in io.Reader) error {
	// initialize lexer and set compiler
	lex := NewLexerWithInit(in, func(y *Lexer) { y.p = c })
	yyErrorVerbose = true
	ret := yyParse(lex)

	if ret == 1 {
		return c.err
	}
	return nil
}

// PlotAst plots AST into file
func (c *Compiler) PlotAst(filename string) error {
	if filename == "" {
		return fmt.Errorf("empty file name")
	}
	err := GenerateDotFormat(c.ast, filename)
	if err != nil {
		return err
	}
	return nil
}
