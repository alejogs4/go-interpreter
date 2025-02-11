package token

type TokenType string

const (
	Illegal = TokenType("ILLEGAL")
	EOF     = TokenType("EOF")

	Ident = TokenType("IDENT")
	INT   = TokenType("INT")

	Assign = TokenType("=")
	Plus   = TokenType("+")
	Minus  = TokenType("-")

	Comma     = TokenType(",")
	Semicolon = TokenType(";")

	LParen  = TokenType("(")
	RParent = TokenType(")")
	LBrace  = TokenType("{")
	RBrace  = TokenType("}")

	Function = TokenType("FUNCTION")
	Let      = TokenType("LET")

	Negation   = TokenType("!")
	Slash      = TokenType("/")
	Equal      = TokenType("==")
	Different  = TokenType("!=")
	Product    = TokenType("*")
	LessThan   = TokenType("<")
	BiggerThan = TokenType(">")

	True            = TokenType("true")
	False           = TokenType("false")
	IfConditional   = TokenType("if")
	ElseConditional = TokenType("else")
	ReturnStatement = TokenType("return")
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     IfConditional,
	"else":   ElseConditional,
	"return": ReturnStatement,
}

func LookupIdentifier(identifier string) TokenType {
	if ident, ok := keywords[identifier]; ok {
		return ident
	}

	return Ident
}
