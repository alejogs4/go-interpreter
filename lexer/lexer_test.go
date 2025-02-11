package lexer

import (
	"testing"

	"go-interpreter.com/m/token"
)

func TestNew(t *testing.T) {
	input := `=+(){},;`
	tests := []struct {
		name            string
		expectedToken   token.TokenType
		expectedLiteral string
	}{
		{
			name:            "Assign",
			expectedToken:   token.Assign,
			expectedLiteral: "=",
		},
		{
			name:            "Plus",
			expectedToken:   token.Plus,
			expectedLiteral: "+",
		},
		{
			name:            "Left parenthesis",
			expectedToken:   token.LParen,
			expectedLiteral: "(",
		},
		{
			name:            "Right parenthesis",
			expectedToken:   token.RParent,
			expectedLiteral: ")",
		},
		{
			name:            "Left bracket",
			expectedToken:   token.LBrace,
			expectedLiteral: "{",
		},
		{
			name:            "Right bracket",
			expectedToken:   token.RBrace,
			expectedLiteral: "}",
		},
		{
			name:            "Comma",
			expectedToken:   token.Comma,
			expectedLiteral: ",",
		},
		{
			name:            "Semicolon",
			expectedToken:   token.Semicolon,
			expectedLiteral: ";",
		},
		{
			name:            "EOF",
			expectedToken:   token.EOF,
			expectedLiteral: "\x00",
		},
	}

	lexer := New(input)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken := lexer.NextToken()
			if gotToken.Type != tt.expectedToken {
				t.Fatalf("Token type wrong expected=%q actual=%q", tt.expectedToken, gotToken.Type)
			}

			if gotToken.Literal != tt.expectedLiteral {
				t.Fatalf("Token literal wrong expected=%q actual=%q", tt.expectedLiteral, gotToken.Literal)
			}
		})
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`

	testCases := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.INT, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.INT, "10"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParent, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RParent, ")"},
		{token.Semicolon, ";"},
		{token.Negation, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Product, "*"},
		{token.INT, "5"},
		{token.Semicolon, ";"},
		{token.INT, "5"},
		{token.LessThan, "<"},
		{token.INT, "10"},
		{token.BiggerThan, ">"},
		{token.INT, "5"},
		{token.Semicolon, ";"},
		{token.IfConditional, "if"},
		{token.LParen, "("},
		{token.INT, "5"},
		{token.LessThan, "<"},
		{token.INT, "10"},
		{token.RParent, ")"},
		{token.LBrace, "{"},
		{token.ReturnStatement, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.ElseConditional, "else"},
		{token.LBrace, "{"},
		{token.ReturnStatement, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.INT, "10"},
		{token.Equal, "=="},
		{token.INT, "10"},
		{token.Semicolon, ";"},
		{token.INT, "10"},
		{token.Different, "!="},
		{token.INT, "9"},
		{token.Semicolon, ";"},
		{token.EOF, "\x00"},
	}

	lexer := New(input)
	t.Run("NextToken test", func(t *testing.T) {
		for _, tt := range testCases {
			gotToken := lexer.NextToken()
			if gotToken.Type != tt.expectedType {
				t.Fatalf("Token type wrong expected=%q actual=%q", tt.expectedType, gotToken.Type)
			}

			if gotToken.Literal != tt.expectedLiteral {
				t.Fatalf("Token literal wrong expected=%q actual=%q", tt.expectedLiteral, gotToken.Literal)
			}
		}
	})
}
