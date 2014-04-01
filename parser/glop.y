%{

package parser

import (
  "github.com/palats/glop/nodes"
)

%}

%union {
  raw string
  node nodes.Node
}

%token <raw> tokOpen
%token <raw> tokClose
%token <raw> tokIdentifier

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
    $$ = nodes.NewExpr($1, nil)
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

%%

