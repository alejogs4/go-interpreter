package lexer

import "go-interpreter.com/m/token"

type Lexer struct {
	code         string
	position     int  // Current char
	readPosition int  // Next char
	ch           byte // Char under examination
}

func New(sourceCode string) *Lexer {
	l := &Lexer{code: sourceCode}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			currentCharacter := l.ch
			l.readChar()
			literal := string(currentCharacter) + string(l.ch)
			tok = token.Token{Type: token.Equal, Literal: literal}
		} else {
			tok = newToken(token.Assign, l.ch)
		}
	case ';':
		tok = newToken(token.Semicolon, l.ch)
	case '(':
		tok = newToken(token.LParen, l.ch)
	case ')':
		tok = newToken(token.RParent, l.ch)
	case '!':
		if l.peekChar() == '=' {
			currentCharacter := l.ch
			l.readChar()
			literal := string(currentCharacter) + string(l.ch)
			tok = token.Token{Type: token.Different, Literal: literal}
		} else {
			tok = newToken(token.Negation, l.ch)
		}
	case '-':
		tok = newToken(token.Minus, l.ch)
	case '/':
		tok = newToken(token.Slash, l.ch)
	case '*':
		tok = newToken(token.Product, l.ch)
	case '<':
		tok = newToken(token.LessThan, l.ch)
	case '>':
		tok = newToken(token.BiggerThan, l.ch)
	case ',':
		tok = newToken(token.Comma, l.ch)
	case '+':
		tok = newToken(token.Plus, l.ch)
	case '{':
		tok = newToken(token.LBrace, l.ch)
	case '}':
		tok = newToken(token.RBrace, l.ch)
	case 0:
		tok = newToken(token.EOF, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readDigit()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.Illegal, l.ch)
		}

	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.code) {
		return 0
	} else {
		return l.code[l.readPosition]
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.code) {
		l.ch = 0
	} else {
		l.ch = l.code[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	initialPosition := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.code[initialPosition:l.position]
}

func (l *Lexer) readDigit() string {
	initialPosition := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.code[initialPosition:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return ch <= 'z' && ch >= 'a' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
