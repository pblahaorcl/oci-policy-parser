package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var file, statement string
	flag.StringVar(&file, "file", "output.dot", "DOT File name ")
	flag.StringVar(&statement, "statement", "Allow group blaha to use vcns in bar where foo == 2", "statement")
	flag.Parse()

	comp := NewCompiler()
	log.Printf("parsing: %s in dot file: %s", statement, file)
	err := comp.Parse(strings.NewReader(statement))
	if err != nil {
		switch e := err.(type) {
		case *SyntaxError:
			fmt.Printf("%s at column: %d \n", e.Msg, e.Col)
			os.Exit(127)
		default:
			fmt.Println(e)
		}
	}
	err = comp.PlotAst(file)
	if err != nil {
		fmt.Printf("%s", err)
	}
}
