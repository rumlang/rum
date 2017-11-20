package parser

// Context is provided to the token function when the topdown parser is
// proceeding through the data.
type Context interface {
	// Expression allows a token nud/led function to call a sub-expression
	// parsing (e.g., for things within parenthesis).
	Expression(rbp int) interface{}
	// Advance tells the parse to discard the incoming token and move to the next
	// one.
	Advance() Token
	// Peek returns the incoming token. That can be used to check whether the
	// next token is what is expected (and then discard it with Advance()).
	Peek() Token
	// Error indicates that something went wrong. The parsing will still continue
	// after that, allowing for multiple errors to be found.
	Error(error)
}

// Token must be implemented by the data returned by the lexer provided to the
// parser.
type Token interface {
	// Nud is 'Null denotation'. This is used in a top down parser when a token
	// is encountered at the beginning of an expression.
	Nud(ctx Context) interface{}
	// Led is 'Left denotation'. This is used in a top down parser when a token
	// is encountered after the beginning of the expression. 'left' contains what
	// was previously obtained with Nud/Led of the previous token.
	Led(ctx Context, left interface{}) interface{}
	// Lbp is 'Left Binding Priority'. The value 0 is used as right binding
	// priority in TopDownParse - which usually means that "end of stream" should
	// have an lbp of 0.
	Lbp() int
}

// Lexer is the interface that the parser expects to get input tokens.
type Lexer interface {
	// Next provide the next available token.
	Next() Token
}

// TopDown implements a generic TopDown parser.
type TopDown struct {
	lex Lexer

	// token contains the next token to be encountered.
	token Token

	// All errors coming from the tokens nud/led functions.
	errors []error
}

// Advance returns the incoming token, and move to the next one.
func (p *TopDown) Advance() Token {
	t := p.token
	p.token = p.lex.Next()
	return t
}

// Peek returns the incoming token.
func (p *TopDown) Peek() Token {
	return p.token
}

// Expression continue the parsing and returns the parsed data. The actual data
// depends on what nud&led functions of the tokens coming from the lexer
// return.
// Parameter rbp is the right-binding-priority; when parsing a full expression,
// it should be zero usually (it depends on the priorities the tokens provide;
// usually, a token with priority zero indicates end of stream).
func (p *TopDown) Expression(rbp int) interface{} {
	left := p.Advance().Nud(p)
	for rbp < p.Peek().Lbp() {
		left = p.Advance().Led(p, left)
	}
	return left
}

// Error implements the Context interface and allows a token nud&led to
// indicates that there was an issue.
func (p *TopDown) Error(err error) {
	p.errors = append(p.errors, err)
}

// TopDownParse run the topdown parser on the tokens provided by the lexer.
// It will return whatever the tokens have built through nud/led methods and
// the list of errors encountered during the parsing.
func TopDownParse(lex Lexer) (interface{}, []error) {
	p := &TopDown{
		lex: lex,
	}
	// Make the first token available.
	p.Advance()
	return p.Expression(0), p.errors
}
