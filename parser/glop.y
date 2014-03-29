%{

package parser

%}

%union {
  raw string
}

%token <raw> tokOpen
%token <raw> tokClose
%token <raw> tokIdentifier

%type <raw> program expr multiplop plop

%%

program: expr

expr: tokOpen multiplop tokClose

multiplop:
         plop | plop multiplop

plop: tokIdentifier | expr

%%

