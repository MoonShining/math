package lexer

type TokenType int

const (
	Add TokenType = iota
	Sub
	Mul
	Div
	LeftParen
	RightParen
	Number
	Eof
)

var (
	EOF         = Token{Type: Eof}
	ADD         = Token{Type: Add}
	SUB         = Token{Type: Sub}
	MUL         = Token{Type: Mul}
	DIV         = Token{Type: Div}
	LEFT_PAREN  = Token{Type: LeftParen}
	RIGHT_PAREN = Token{Type: RightParen}
)

func NewLexer(in string) *Lexer {
	l := &Lexer{input: []rune(in)}
	l.readChar()
	return l
}

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input []rune

	position     int
	readPosition int
	ch           rune
}

func (l *Lexer) NextToken() Token {
	switch l.ch {
	case '-':
		l.readChar()
		return SUB
	case '+':
		l.readChar()
		return ADD
	case '*':
		l.readChar()
		return MUL
	case '/':
		l.readChar()
		return DIV
	case '(':
		l.readChar()
		return LEFT_PAREN
	case ')':
		l.readChar()
		return RIGHT_PAREN
	case 0:
		return EOF
	default:
		literal := []rune{}
		for l.isDigit() {
			literal = append(literal, l.ch)
			l.readChar()
		}
		return Token{Type: Number, Literal: string(literal)}
	}
}

func (l *Lexer) isDigit() bool {
	return '0' <= l.ch && l.ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}
