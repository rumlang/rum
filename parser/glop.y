%{

package parser

import (
  "github.com/palats/glop/nodes"
)

%}

%union {
  token tokenInfo
  node nodes.Node
}

%token <token> tokOpen
%token <token> tokClose
%token <token> tokIdentifier
%token <token> tokInteger

%type <node> program expr atom list

%%

program: expr
  {
    yylex.(*lexer).program = $1
  }


expr:
  tokOpen list tokClose
  {
    $$ = $2
  }

| tokOpen tokClose
  {
    $$ = nodes.NewExpr(nil, nil)
  }

| atom
  {
    $$ = $1
  }


list:
  expr
  {
    $$ = nodes.NewExpr($1, nil)
  }
| expr list
  {
    $$ = nodes.NewExpr($1, $2)
  }


atom: 
  tokIdentifier
  {
    $$ = nodes.NewIdentifier($1)
  }
| tokInteger
  {
    $$ = nodes.NewInteger($1)
  }

%%

