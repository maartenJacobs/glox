package internal

import (
	"fmt"
	"strconv"
)

// Define all token types.
const (
	// TODO: hardcode values instead of relying on iota. That will allow future compatibility
	// as the value should not change over time, e.g. by entering a token between 2 tokens.

	// Single-character tokens.
	TokenLeftParen = iota
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemicolon
	TokenSlash
	TokenStar

	// One or two character tokens.
	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual

	// Literals.
	TokenIdentifier
	TokenString // Our strings are multiline. Preceding spaces are not trimmed.
	TokenNumber

	// Keywords.
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFun
	TokenFor
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile

	TokenEof
)

// Token represents a lexeme read from the input code, the inferred type and the location
// in the input code. The literal value is also included, if any is interpreted.
type Token struct {
	Type    int // One of the TOKEN_* constants
	Lexeme  string
	Literal interface{}
	Line    int
}

func (token Token) String() string {
	// Horrible translation but more readable than the integer.
	tokenType := fmt.Sprintf("%d", token.Type)
	if token.Type == TokenAnd {
		tokenType = "AND"
	} else if token.Type == TokenLeftParen {
		tokenType = "LEFT_PAREN"
	} else if token.Type == TokenRightParen {
		tokenType = "RIGHT_PAREN"
	} else if token.Type == TokenLeftBrace {
		tokenType = "LEFT_BRACE"
	} else if token.Type == TokenRightBrace {
		tokenType = "RIGHT_BRACE"
	} else if token.Type == TokenComma {
		tokenType = "COMMA"
	} else if token.Type == TokenDot {
		tokenType = "DOT"
	} else if token.Type == TokenMinus {
		tokenType = "MINUS"
	} else if token.Type == TokenPlus {
		tokenType = "PLUS"
	} else if token.Type == TokenSemicolon {
		tokenType = "SEMICOLON"
	} else if token.Type == TokenSlash {
		tokenType = "SLASH"
	} else if token.Type == TokenStar {
		tokenType = "STAR"
	} else if token.Type == TokenBang {
		tokenType = "BANG"
	} else if token.Type == TokenBangEqual {
		tokenType = "BANG_EQUAL"
	} else if token.Type == TokenEqual {
		tokenType = "EQUAL"
	} else if token.Type == TokenEqualEqual {
		tokenType = "EQUAL_EQUAL"
	} else if token.Type == TokenGreater {
		tokenType = "GREATER"
	} else if token.Type == TokenGreaterEqual {
		tokenType = "GREATER_EQUAL"
	} else if token.Type == TokenLess {
		tokenType = "LESS"
	} else if token.Type == TokenLessEqual {
		tokenType = "LESS_EQUAL"
	} else if token.Type == TokenIdentifier {
		tokenType = "IDENTIFIER"
	} else if token.Type == TokenString {
		tokenType = "STRING"
	} else if token.Type == TokenNumber {
		tokenType = "NUMBER"
	} else if token.Type == TokenAnd {
		tokenType = "AND"
	} else if token.Type == TokenClass {
		tokenType = "CLASS"
	} else if token.Type == TokenElse {
		tokenType = "ELSE"
	} else if token.Type == TokenFalse {
		tokenType = "FALSE"
	} else if token.Type == TokenFun {
		tokenType = "FUN"
	} else if token.Type == TokenFor {
		tokenType = "FOR"
	} else if token.Type == TokenIf {
		tokenType = "IF"
	} else if token.Type == TokenNil {
		tokenType = "NIL"
	} else if token.Type == TokenOr {
		tokenType = "OR"
	} else if token.Type == TokenPrint {
		tokenType = "PRINT"
	} else if token.Type == TokenReturn {
		tokenType = "RETURN"
	} else if token.Type == TokenSuper {
		tokenType = "SUPER"
	} else if token.Type == TokenThis {
		tokenType = "THIS"
	} else if token.Type == TokenTrue {
		tokenType = "TRUE"
	} else if token.Type == TokenVar {
		tokenType = "VAR"
	} else if token.Type == TokenWhile {
		tokenType = "WHILE"
	} else if token.Type == TokenEof {
		tokenType = "EOF"
	}

	return fmt.Sprintf("%s %s %v", tokenType, token.Lexeme, token.Literal)
}

// Define all the keywords
var keywords = map[string]int{
	"and": TokenAnd,
	"class": TokenClass,
	"else": TokenElse,
	"false": TokenFalse,
	"for": TokenFor,
	"fun": TokenFun,
	"if": TokenIf,
	"nil": TokenNil,
	"or": TokenOr,
	"print": TokenPrint,
	"return": TokenReturn,
	"super": TokenSuper,
	"this": TokenThis,
	"true": TokenTrue,
	"var": TokenVar,
	"while": TokenWhile,
}

// Scanner scans the source code left to right and returns a list of tokens interpreted from
// the source code.
type Scanner struct {
	source   []byte
	reporter ErrorReporter
	// Scanning state:
	start   int     // The location of the first character in the current lexeme being scanned
	current int     // The location of the current character in the current lexeme being scanned
	line    int     // The line number of the current position in the code
	tokens  []Token // Scanned tokens
}

func NewScanner(source []byte, reporter ErrorReporter) Scanner {
	return Scanner{
		source:   source,
		reporter: reporter,
		start:    0,
		current:  0,
		line:     1,
	}
}

func (scanner *Scanner) ScanTokens() []Token {
	for !scanner.isAtEnd() {
		// We are at the beginning of the next lexeme.
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, Token{TokenEof, "", nil, scanner.line})
	return scanner.tokens
}

func (scanner Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

func (scanner *Scanner) scanToken() {
	switch c := scanner.advance(); c {
	case '(':
		scanner.addToken(TokenLeftParen)
	case ')':
		scanner.addToken(TokenRightParen)
	case '{':
		scanner.addToken(TokenLeftBrace)
	case '}':
		scanner.addToken(TokenRightBrace)
	case ',':
		scanner.addToken(TokenComma)
	case '.':
		scanner.addToken(TokenDot)
	case '-':
		scanner.addToken(TokenMinus)
	case '+':
		scanner.addToken(TokenPlus)
	case ';':
		scanner.addToken(TokenSemicolon)
	case '*':
		scanner.addToken(TokenStar)
	case '!':
		if scanner.match('=') {
			scanner.addToken(TokenBangEqual)
		} else {
			scanner.addToken(TokenBang)
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(TokenEqualEqual)
		} else {
			scanner.addToken(TokenEqual)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(TokenLessEqual)
		} else {
			scanner.addToken(TokenLess)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(TokenGreaterEqual)
		} else {
			scanner.addToken(TokenGreater)
		}
	case '/':
		if scanner.match('/') {
			// We use `peek` here to keep the newline under consideration. We'll advance
			// the line number counter if there is one but that's a different part of the loop.
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else {
			scanner.addToken(TokenSlash)
		}
	case ' ':
	case '\t':
	case '\r':
		// Ignore whitespace.
		break
	case '\n':
		scanner.line++
	case '"':
		scanner.string()
	default:
		if scanner.isDigit(c) {
			scanner.number()
		} else if scanner.isAlpha(c) {
			scanner.identifier()
		} else {
			scanner.reporter.Error(scanner.line, "Unexpected character.")
		}
	}
}

func (scanner *Scanner) advance() byte {
	scanner.current++
	return scanner.source[scanner.current-1]
}

func (scanner *Scanner) addToken(tokenType int) {
	scanner.addLiteralToken(tokenType, nil)
}

func (scanner *Scanner) addLiteralToken(tokenType int, literal interface{}) {
	text := string(scanner.source[scanner.start:scanner.current])
	scanner.tokens = append(scanner.tokens, Token{tokenType, text, literal, scanner.line})
}

// Match is a conditional advance.
func (scanner *Scanner) match(expected byte) bool {
	if scanner.isAtEnd() {
		return false
	}
	if scanner.source[scanner.current] != expected {
		return false
	}

	scanner.current++
	return true
}

func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	}
	return scanner.source[scanner.current]
}

func (scanner *Scanner) string() {
	// Scan until string or input end.
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}

	// Unterminated string.
	if scanner.isAtEnd() {
		scanner.reporter.Error(scanner.line, "Unterminated string.")
		return
	}

	// The closing ".
	scanner.advance()

	// Trim the surrounding quotes to track the string token.
	value := string(scanner.source[scanner.start+1 : scanner.current-1])
	scanner.addLiteralToken(TokenString, value)
}

func (scanner *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (scanner *Scanner) number() {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && scanner.isDigit(scanner.peekNext()) {
		scanner.advance()
		for scanner.isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	value := string(scanner.source[scanner.start:scanner.current])
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil { // This would be due to a compiler programmer's error
		panic(err)
	}
	scanner.addLiteralToken(TokenNumber, floatValue)
}

func (scanner *Scanner) peekNext() byte {
	if scanner.current+1 >= len(scanner.source) {
		return 0
	}
	return scanner.source[scanner.current+1]
}

func (scanner Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (scanner Scanner) isAlphaNumeric(c byte) bool {
	return scanner.isAlpha(c) || scanner.isDigit(c)
}

func (scanner *Scanner) identifier() {
	for scanner.isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	// See if the identifier is a reserved word.
	text := string(scanner.source[scanner.start:scanner.current])
	tokenType, found := keywords[text]
	if !found {
		tokenType = TokenIdentifier
	}
	scanner.addToken(tokenType)
}
