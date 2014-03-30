%{

package parser

%}

%union {
  raw string
}

%token <raw> tokOpen
%token <raw> tokClose
%token <raw> tokIdentifier

%type <raw> program expr plop

%%

program: expr

expr:
   plop | plop expr

plop: tokIdentifier | tokOpen tokClose | tokOpen expr tokClose

%%

