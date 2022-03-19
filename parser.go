// Code generated by goyacc -o parser.go parser.y. DO NOT EDIT.

//line parser.y:2
package main

import __yyfmt__ "fmt"

//line parser.y:2

//line parser.y:6
type yySymType struct {
	yys  int
	s    string /* String */
	node Node   /* Node in the AST */
}

const ALLOW = 57346
const TO = 57347
const IN = 57348
const WHERE = 57349
const VERB = 57350
const RESOURCE = 57351
const SUBJECT = 57352
const EQ = 57353
const NE = 57354
const VALUE = 57355
const VARIABLE = 57356

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"ALLOW",
	"TO",
	"IN",
	"WHERE",
	"VERB",
	"RESOURCE",
	"SUBJECT",
	"EQ",
	"NE",
	"VALUE",
	"VARIABLE",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:101

func cast(y yyLexer) *Compiler { return y.(*Lexer).p }

//line yacctab:1
var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 22

var yyAct = [...]int8{
	6, 7, 22, 19, 20, 4, 11, 9, 15, 12,
	5, 2, 21, 14, 18, 13, 17, 10, 8, 3,
	16, 1,
}

var yyPact = [...]int16{
	7, -1000, -5, 5, -13, -1, -1000, -1000, -3, -1000,
	3, -1000, -13, 1, -1000, -13, -1000, -8, -11, -1000,
	-1000, -1000, -1000,
}

var yyPgo = [...]int8{
	0, 21, 20, 19, 18, 17, 15, 14, 0, 12,
}

var yyR1 = [...]int8{
	0, 1, 1, 6, 5, 4, 3, 8, 9, 7,
	7, 2,
}

var yyR2 = [...]int8{
	0, 9, 7, 1, 1, 1, 2, 1, 1, 1,
	1, 3,
}

var yyChk = [...]int16{
	-1000, -1, 4, -3, 10, 5, -8, 14, -4, 8,
	-5, 9, 6, -6, -8, 7, -2, -8, -7, 11,
	12, -9, 13,
}

var yyDef = [...]int8{
	0, -2, 0, 0, 0, 0, 6, 7, 0, 5,
	0, 4, 0, 2, 3, 0, 1, 0, 0, 9,
	10, 11, 8,
}

var yyTok1 = [...]int8{
	1,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14,
}

var yyTok3 = [...]int8{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(yyPact[state])
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && int(yyChk[int(yyAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || int(yyExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := int(yyExca[i])
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(yyTok1[0])
		goto out
	}
	if char < len(yyTok1) {
		token = int(yyTok1[char])
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = int(yyTok2[char-yyPrivate])
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = int(yyTok3[i+0])
		if token == char {
			token = int(yyTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(yyTok2[1]) /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = int(yyPact[yystate])
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = int(yyAct[yyn])
	if int(yyChk[yyn]) == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = int(yyDef[yystate])
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && int(yyExca[xi+1]) == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = int(yyExca[xi+0])
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = int(yyExca[xi+1])
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = int(yyPact[yyS[yyp].yys]) + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = int(yyAct[yyn]) /* simulate a shift of "error" */
					if int(yyChk[yystate]) == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= int(yyR2[yyn])
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = int(yyR1[yyn])
	yyg := int(yyPgo[yyn])
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = int(yyAct[yyg])
	} else {
		yystate = int(yyAct[yyj])
		if int(yyChk[yystate]) != -yyn {
			yystate = int(yyAct[yyg])
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.y:29
		{
			yyVAL.node = newStatementNode(yyDollar[2].node, yyDollar[4].node, yyDollar[5].node, yyDollar[7].node, yyDollar[9].node)
			cast(yylex).SetAstRoot(yyVAL.node)
		}
	case 2:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.y:34
		{
			yyVAL.node = newStatementNode(yyDollar[2].node, yyDollar[4].node, yyDollar[5].node, yyDollar[7].node, nil)
			cast(yylex).SetAstRoot(yyVAL.node)
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:42
		{
			yyVAL.node = newTokenNode(NODE_LOCATION, "", yyDollar[1].node)
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:49
		{
			yyVAL.node = newTokenNode(NODE_RESOURCE, yyDollar[1].s, nil)
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:56
		{
			yyVAL.node = newTokenNode(NODE_VERB, yyDollar[1].s, nil)
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.y:63
		{
			yyVAL.node = newTokenNode(NODE_SUBJECT, "", yyDollar[2].node)
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:70
		{
			yyVAL.node = newTokenNode(NODE_VARIABLE, yyDollar[1].s, nil)
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:77
		{
			yyVAL.node = newTokenNode(NODE_VALUE, yyDollar[1].s, nil)
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:84
		{
			yyVAL.node = newTokenNode(NODE_OP, yyDollar[1].s, nil)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.y:88
		{
			yyVAL.node = newTokenNode(NODE_OP, yyDollar[1].s, nil)
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.y:95
		{
			yyVAL.node = newConditionNode(yyDollar[1].node, yyDollar[2].node, yyDollar[3].node)
		}
	}
	goto yystack /* stack new state and value */
}
