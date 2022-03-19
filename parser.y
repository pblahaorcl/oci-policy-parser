%{
package main

%}

%union {
    s string    /* String */
    node Node   /* Node in the AST */
};


%token ALLOW
%token TO
%token IN
%token WHERE
%token VERB
%token RESOURCE
%token SUBJECT
%token EQ
%token NE

%type <node> statement condition subj verb resource location op v value
%token <s> VALUE VARIABLE VERB RESOURCE EQ NE

%%

statement:
    ALLOW subj TO verb resource IN location WHERE condition
    {
        $$ = newTokenNode(NODE_STATEMENT, "", $2, $4, $5, $7, $9)
        cast(yylex).SetAstRoot($$)
    }
    | ALLOW subj TO verb resource IN location
    {
        $$ = newTokenNode(NODE_STATEMENT, "", $2, $4, $5, $7, nil)
        cast(yylex).SetAstRoot($$)
    }
    ;

location:
    v
    {
        $$ = newTokenNode(NODE_LOCATION, "", $1)
    }
    ;

resource:
    RESOURCE
    {
        $$ = newTokenNode(NODE_RESOURCE, $1, nil)
    }
    ;

verb:
    VERB
    {
        $$ = newTokenNode(NODE_VERB, $1, nil)
    }
    ;

subj:
    SUBJECT v
    {
        $$ = newTokenNode(NODE_SUBJECT, "", $2)
    }
    ;

v:
    VARIABLE
    {
        $$ = newTokenNode(NODE_VARIABLE, $1, nil)
    }
    ;

value:
    VALUE
    {
        $$ = newTokenNode(NODE_VALUE, $1, nil)
    }
    ;

op:
    EQ
    {
        $$ = newTokenNode(NODE_OP, $1, nil)
    }
    | NE
    {
        $$ = newTokenNode(NODE_OP, $1, nil)
    }
    ;

condition:
    v op value
    {
        $$ = newTokenNode(NODE_CONDITION, "", $1, $2, $3)
    }
    ;


%%
func cast(y yyLexer) *Compiler { return y.(*Lexer).p }
