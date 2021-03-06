package main

import "fmt"
import (
	"bufio"
	"io"
	"strings"
)

type frame struct {
	i            int
	s            string
	line, column int
}
type Lexer struct {
	// The lexer runs in its own goroutine, and communicates via channel 'ch'.
	ch      chan frame
	ch_stop chan bool
	// We record the level of nesting because the action could return, and a
	// subsequent call expects to pick up where it left off. In other words,
	// we're simulating a coroutine.
	// TODO: Support a channel-based variant that compatible with Go's yacc.
	stack []frame
	stale bool

	// The 'l' and 'c' fields were added for
	// https://github.com/wagerlabs/docker/blob/65694e801a7b80930961d70c69cba9f2465459be/buildfile.nex
	// Since then, I introduced the built-in Line() and Column() functions.
	l, c int

	parseResult interface{}

	// The following line makes it easy for scripts to insert fields in the
	// generated code.
p *Compiler
	// [NEX_END_OF_LEXER_STRUCT]
}

// NewLexerWithInit creates a new Lexer object, runs the given callback on it,
// then returns it.
func NewLexerWithInit(in io.Reader, initFun func(*Lexer)) *Lexer {
	yylex := new(Lexer)
	if initFun != nil {
		initFun(yylex)
	}
	yylex.ch = make(chan frame)
	yylex.ch_stop = make(chan bool, 1)
	var scan func(in *bufio.Reader, ch chan frame, ch_stop chan bool, family []dfa, line, column int)
	scan = func(in *bufio.Reader, ch chan frame, ch_stop chan bool, family []dfa, line, column int) {
		// Index of DFA and length of highest-precedence match so far.
		matchi, matchn := 0, -1
		var buf []rune
		n := 0
		checkAccept := func(i int, st int) bool {
			// Higher precedence match? DFAs are run in parallel, so matchn is at most len(buf), hence we may omit the length equality check.
			if family[i].acc[st] && (matchn < n || matchi > i) {
				matchi, matchn = i, n
				return true
			}
			return false
		}
		var state [][2]int
		for i := 0; i < len(family); i++ {
			mark := make([]bool, len(family[i].startf))
			// Every DFA starts at state 0.
			st := 0
			for {
				state = append(state, [2]int{i, st})
				mark[st] = true
				// As we're at the start of input, follow all ^ transitions and append to our list of start states.
				st = family[i].startf[st]
				if -1 == st || mark[st] {
					break
				}
				// We only check for a match after at least one transition.
				checkAccept(i, st)
			}
		}
		atEOF := false
		stopped := false
		for {
			if n == len(buf) && !atEOF {
				r, _, err := in.ReadRune()
				switch err {
				case io.EOF:
					atEOF = true
				case nil:
					buf = append(buf, r)
				default:
					panic(err)
				}
			}
			if !atEOF {
				r := buf[n]
				n++
				var nextState [][2]int
				for _, x := range state {
					x[1] = family[x[0]].f[x[1]](r)
					if -1 == x[1] {
						continue
					}
					nextState = append(nextState, x)
					checkAccept(x[0], x[1])
				}
				state = nextState
			} else {
			dollar: // Handle $.
				for _, x := range state {
					mark := make([]bool, len(family[x[0]].endf))
					for {
						mark[x[1]] = true
						x[1] = family[x[0]].endf[x[1]]
						if -1 == x[1] || mark[x[1]] {
							break
						}
						if checkAccept(x[0], x[1]) {
							// Unlike before, we can break off the search. Now that we're at the end, there's no need to maintain the state of each DFA.
							break dollar
						}
					}
				}
				state = nil
			}

			if state == nil {
				lcUpdate := func(r rune) {
					if r == '\n' {
						line++
						column = 0
					} else {
						column++
					}
				}
				// All DFAs stuck. Return last match if it exists, otherwise advance by one rune and restart all DFAs.
				if matchn == -1 {
					if len(buf) == 0 { // This can only happen at the end of input.
						break
					}
					lcUpdate(buf[0])
					buf = buf[1:]
				} else {
					text := string(buf[:matchn])
					buf = buf[matchn:]
					matchn = -1
					select {
					case ch <- frame{matchi, text, line, column}:
						{
						}
					case stopped = <-ch_stop:
						{
						}
					}
					if stopped {
						break
					}
					if len(family[matchi].nest) > 0 {
						scan(bufio.NewReader(strings.NewReader(text)), ch, ch_stop, family[matchi].nest, line, column)
					}
					if atEOF {
						break
					}
					for _, r := range text {
						lcUpdate(r)
					}
				}
				n = 0
				for i := 0; i < len(family); i++ {
					state = append(state, [2]int{i, 0})
				}
			}
		}
		ch <- frame{-1, "", line, column}
	}
	go scan(bufio.NewReader(in), yylex.ch, yylex.ch_stop, dfas, 0, 0)
	return yylex
}

type dfa struct {
	acc          []bool           // Accepting states.
	f            []func(rune) int // Transitions.
	startf, endf []int            // Transitions at start and end of input.
	nest         []dfa
}

var dfas = []dfa{
	// Allow|ALLOW|allow
	{[]bool{false, false, false, false, false, false, true, false, false, false, false, true, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 65:
				return 1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return 2
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return 7
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return 8
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return 3
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return 4
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return 5
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return 6
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return 12
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return 9
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return 10
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return 13
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return 14
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 65:
				return -1
			case 76:
				return -1
			case 79:
				return -1
			case 87:
				return -1
			case 97:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			case 119:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// to|TO
	{[]bool{false, false, false, true, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 79:
				return -1
			case 84:
				return 1
			case 111:
				return -1
			case 116:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 79:
				return 4
			case 84:
				return -1
			case 111:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 79:
				return -1
			case 84:
				return -1
			case 111:
				return 3
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 79:
				return -1
			case 84:
				return -1
			case 111:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 79:
				return -1
			case 84:
				return -1
			case 111:
				return -1
			case 116:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// in
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 105:
				return 1
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 105:
				return -1
			case 110:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 105:
				return -1
			case 110:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// where
	{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 104:
				return -1
			case 114:
				return -1
			case 119:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 104:
				return 2
			case 114:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return 3
			case 104:
				return -1
			case 114:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 104:
				return -1
			case 114:
				return 4
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return 5
			case 104:
				return -1
			case 114:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 104:
				return -1
			case 114:
				return -1
			case 119:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

	// manage|use|read
	{[]bool{false, false, false, false, false, true, false, false, true, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return 1
			case 110:
				return -1
			case 114:
				return 2
			case 115:
				return -1
			case 117:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return 9
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return 6
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return 4
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return 5
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return 7
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return 8
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return 10
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return 11
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return 12
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return 13
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 103:
				return -1
			case 109:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 117:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// vcns|subnets
	{[]bool{false, false, false, false, false, true, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return 1
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return 6
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return 3
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return 4
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return 5
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return 7
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return 8
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return 9
			case 110:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return -1
			case 116:
				return 10
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return 11
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 118:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// group
	{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 103:
				return 1
			case 111:
				return -1
			case 112:
				return -1
			case 114:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 111:
				return -1
			case 112:
				return -1
			case 114:
				return 2
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 111:
				return 3
			case 112:
				return -1
			case 114:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 111:
				return -1
			case 112:
				return -1
			case 114:
				return -1
			case 117:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 111:
				return -1
			case 112:
				return 5
			case 114:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 111:
				return -1
			case 112:
				return -1
			case 114:
				return -1
			case 117:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

	// [0-9]*
	{[]bool{true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch {
			case 48 <= r && r <= 57:
				return 0
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1}, []int{ /* End-of-input transitions */ -1}, nil},

	// [a-z][a-z0-9]*
	{[]bool{false, true, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch {
			case 48 <= r && r <= 57:
				return -1
			case 97 <= r && r <= 122:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch {
			case 48 <= r && r <= 57:
				return 2
			case 97 <= r && r <= 122:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch {
			case 48 <= r && r <= 57:
				return 2
			case 97 <= r && r <= 122:
				return 2
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// ==
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// !=
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 33:
				return 1
			case 61:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 61:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 61:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// [ \t\n]+
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 9:
				return 1
			case 10:
				return 1
			case 32:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 9:
				return 1
			case 10:
				return 1
			case 32:
				return 1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// .
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			return 1
		},
		func(r rune) int {
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},
}

func NewLexer(in io.Reader) *Lexer {
	return NewLexerWithInit(in, nil)
}

func (yyLex *Lexer) Stop() {
	yyLex.ch_stop <- true
}

// Text returns the matched text.
func (yylex *Lexer) Text() string {
	return yylex.stack[len(yylex.stack)-1].s
}

// Line returns the current line number.
// The first line is 0.
func (yylex *Lexer) Line() int {
	if len(yylex.stack) == 0 {
		return 0
	}
	return yylex.stack[len(yylex.stack)-1].line
}

// Column returns the current column number.
// The first column is 0.
func (yylex *Lexer) Column() int {
	if len(yylex.stack) == 0 {
		return 0
	}
	return yylex.stack[len(yylex.stack)-1].column
}

func (yylex *Lexer) next(lvl int) int {
	if lvl == len(yylex.stack) {
		l, c := 0, 0
		if lvl > 0 {
			l, c = yylex.stack[lvl-1].line, yylex.stack[lvl-1].column
		}
		yylex.stack = append(yylex.stack, frame{0, "", l, c})
	}
	if lvl == len(yylex.stack)-1 {
		p := &yylex.stack[lvl]
		*p = <-yylex.ch
		yylex.stale = false
	} else {
		yylex.stale = true
	}
	return yylex.stack[lvl].i
}
func (yylex *Lexer) pop() {
	yylex.stack = yylex.stack[:len(yylex.stack)-1]
}

// Lex runs the lexer. Always returns 0.
// When the -s option is given, this function is not generated;
// instead, the NN_FUN macro runs the lexer.
func (yylex *Lexer) Lex(lval *yySymType) int {
OUTER0:
	for {
		switch yylex.next(0) {
		case 0:
			{
				return ALLOW
			}
		case 1:
			{
				return TO
			}
		case 2:
			{
				return IN
			}
		case 3:
			{
				return WHERE
			}
		case 4:
			{
				lval.s = yylex.Text()
				return VERB
			}
		case 5:
			{
				lval.s = yylex.Text()
				return RESOURCE
			}
		case 6:
			{
				return SUBJECT
			}
		case 7:
			{
				lval.s = yylex.Text()
				return VALUE
			}
		case 8:
			{
				lval.s = yylex.Text()
				return VARIABLE
			}
		case 9:
			{
				lval.s = yylex.Text()
				return EQ
			}
		case 10:
			{
				lval.s = yylex.Text()
				return NE
			}
		case 11:
			{ /* eat up whitespace */
			}
		case 12:
			{
				fmt.Sprintf("Unrecognized character: %s", yylex.Text())
			}
		default:
			break OUTER0
		}
		continue
	}
	yylex.pop()

	return 0
}
func (yylex Lexer) Error(e string) {
	yylex.p.err = &SyntaxError{
		Row: yylex.Line() + 1,
		Col: yylex.Column(),
		Msg: e,
	}
}
