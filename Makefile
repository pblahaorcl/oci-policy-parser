BINARY=parser

# Build the project
all: clean build test


setup:
	go install github.com/blynn/nex@latest
	go install golang.org/x/tools/cmd/goyacc@latest

build:
	nex -e=true lexer.nex
	printf '/NEX_END_OF_LEXER_STRUCT/i\np *Compiler\n.\nw\nq\n' | ed -s lexer.nn.go
	goyacc -o parser.go parser.y
	go build -o bin/parser .

test:
	go test -v .

run:
	./bin/parser -statement 'Allow group blaha to use vcns in bar where foo == 2'

ast:
	dot -Tpdf output.dot -o output.pdf && open output.pdf

create_golden_file:
	cp output-test.dot testdata/golden.file
	dot -Tpdf testdata/golden.file -o output.pdf && open output.pdf


fuzzy:
	go test -v -tags=fuzzy -fuzz Fuzz -fuzztime=120s .

clean:
	rm -rf bin lexer.nn.go parser.go output.dot *.pdf


.PHONY: setup build clean run test ast create_golden_file fuzzy
