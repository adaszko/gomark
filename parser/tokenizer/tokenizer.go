// Package tokenizer converts markdown text into lexical tokens used by the parser.
package tokenizer

import "unicode/utf8"

type TokenType = string

// Special character tokens.
const (
	Underscore         TokenType = "_"
	Asterisk           TokenType = "*"
	PoundSign          TokenType = "#"
	Backtick           TokenType = "`"
	LeftSquareBracket  TokenType = "["
	RightSquareBracket TokenType = "]"
	LeftParenthesis    TokenType = "("
	RightParenthesis   TokenType = ")"
	ExclamationMark    TokenType = "!"
	QuestionMark       TokenType = "?"
	Tilde              TokenType = "~"
	Hyphen             TokenType = "-"
	PlusSign           TokenType = "+"
	Dot                TokenType = "."
	LessThan           TokenType = "<"
	GreaterThan        TokenType = ">"
	DollarSign         TokenType = "$"
	EqualSign          TokenType = "="
	Pipe               TokenType = "|"
	Colon              TokenType = ":"
	Caret              TokenType = "^"
	Apostrophe         TokenType = "'"
	Backslash          TokenType = "\\"
	Slash              TokenType = "/"
	NewLine            TokenType = "\n"
	Space              TokenType = " "
)

// Text based tokens.
const (
	Number TokenType = "number"
	Text   TokenType = ""
)

// Multi-character tokens.
const (
	DoubleAsterisk   TokenType = "**"
	TripleAsterisk   TokenType = "***"
	DoubleUnderscore TokenType = "__"
	TripleUnderscore TokenType = "___"
	DoubleTilde      TokenType = "~~"
	TripleBacktick   TokenType = "```"
	DoubleHyphen     TokenType = "--"
	TripleHyphen     TokenType = "---"
	DoubleEqual      TokenType = "=="
	DoubleCaret      TokenType = "^^"
	MultipleSpaces   TokenType = "spaces"
)

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(tp, text string) *Token {
	return &Token{
		Type:  tp,
		Value: text,
	}
}

func Tokenize(text string) []*Token {
	tokens := []*Token{}

	byteIndex := 0
	for {
		c, runeWidth := utf8.DecodeRuneInString(text[byteIndex:])
		if c == utf8.RuneError {
			break
		}

		switch c {
		case '_':
			tokens = append(tokens, NewToken(Underscore, "_"))
		case '*':
			tokens = append(tokens, NewToken(Asterisk, "*"))
		case '#':
			tokens = append(tokens, NewToken(PoundSign, "#"))
		case '`':
			tokens = append(tokens, NewToken(Backtick, "`"))
		case '[':
			tokens = append(tokens, NewToken(LeftSquareBracket, "["))
		case ']':
			tokens = append(tokens, NewToken(RightSquareBracket, "]"))
		case '(':
			tokens = append(tokens, NewToken(LeftParenthesis, "("))
		case ')':
			tokens = append(tokens, NewToken(RightParenthesis, ")"))
		case '!':
			tokens = append(tokens, NewToken(ExclamationMark, "!"))
		case '?':
			tokens = append(tokens, NewToken(QuestionMark, "?"))
		case '~':
			tokens = append(tokens, NewToken(Tilde, "~"))
		case '-':
			tokens = append(tokens, NewToken(Hyphen, "-"))
		case '<':
			tokens = append(tokens, NewToken(LessThan, "<"))
		case '>':
			tokens = append(tokens, NewToken(GreaterThan, ">"))
		case '+':
			tokens = append(tokens, NewToken(PlusSign, "+"))
		case '.':
			tokens = append(tokens, NewToken(Dot, "."))
		case '$':
			tokens = append(tokens, NewToken(DollarSign, "$"))
		case '=':
			tokens = append(tokens, NewToken(EqualSign, "="))
		case '|':
			tokens = append(tokens, NewToken(Pipe, "|"))
		case ':':
			tokens = append(tokens, NewToken(Colon, ":"))
		case '^':
			tokens = append(tokens, NewToken(Caret, "^"))
		case '\'':
			tokens = append(tokens, NewToken(Apostrophe, "'"))
		case '\\':
			tokens = append(tokens, NewToken(Backslash, `\`))
		case '/':
			tokens = append(tokens, NewToken(Slash, "/"))
		case '\n':
			tokens = append(tokens, NewToken(NewLine, "\n"))
		case ' ':
			// Handle consecutive spaces
			if token, bytesCount := handleSpaces(text[byteIndex:]); token != nil {
				tokens = append(tokens, token)
				byteIndex += bytesCount
				continue
			}
			tokens = append(tokens, NewToken(Space, " "))
		default:
			// Handle text and numbers more efficiently
			if token, bytesCount := handleTextOrNumber(text[byteIndex:]); token != nil {
				tokens = append(tokens, token)
				byteIndex += bytesCount
				continue
			}
		}

		byteIndex += runeWidth
	}

	return tokens
}

func (t *Token) String() string {
	return t.Value
}

func Stringify(tokens []*Token) string {
	text := ""
	for _, token := range tokens {
		text += token.String()
	}
	return text
}

func Split(tokens []*Token, delimiter TokenType) [][]*Token {
	if len(tokens) == 0 {
		return [][]*Token{}
	}

	result := make([][]*Token, 0)
	current := make([]*Token, 0)
	for _, token := range tokens {
		if token.Type == delimiter {
			result = append(result, current)
			current = make([]*Token, 0)
		} else {
			current = append(current, token)
		}
	}
	result = append(result, current)
	return result
}

func Find(tokens []*Token, target TokenType) int {
	for i, token := range tokens {
		if token.Type == target {
			return i
		}
	}
	return -1
}

func FindUnescaped(tokens []*Token, target TokenType) int {
	for i, token := range tokens {
		if token.Type == target && (i == 0 || (i > 0 && tokens[i-1].Type != Backslash)) {
			return i
		}
	}
	return -1
}

func GetFirstLine(tokens []*Token) []*Token {
	for i, token := range tokens {
		if token.Type == NewLine {
			return tokens[:i]
		}
	}
	return tokens
}

// tryMultiCharToken attempts to match multi-character token sequences.
// Only handles sequences that don't break existing parsers.
// handleSpaces consolidates consecutive spaces into a single token.
func handleSpaces(text string) (*Token, int) {
	runesCount := 0
	bytesCount := 0
	for {
		c, runeWidth := utf8.DecodeRuneInString(text[bytesCount:])
		if c == utf8.RuneError {
			break
		}
		if c == ' ' {
			runesCount++
			bytesCount += runeWidth
		} else {
			break
		}
	}

	if runesCount > 1 {
		spaces := make([]rune, runesCount)
		for i := range spaces {
			spaces[i] = ' '
		}
		return NewToken(MultipleSpaces, string(spaces)), bytesCount
	}

	return nil, 0
}

// handleTextOrNumber consolidates consecutive text or number characters.
func handleTextOrNumber(text string) (*Token, int) {
	c, runeWidth := utf8.DecodeRuneInString(text)
	isNumber := c >= '0' && c <= '9'
	isText := !isSpecialChar(c) && !isNumber

	if !isText && !isNumber {
		return nil, 0
	}

	bytesCount := runeWidth
	for {
		ch, runeWidth := utf8.DecodeRuneInString(text[bytesCount:])
		if ch == utf8.RuneError {
			break
		}
		if isNumber && (ch >= '0' && ch <= '9') {
			bytesCount += runeWidth
		} else if isText && !isSpecialChar(ch) && (ch < '0' || ch > '9') {
			bytesCount += runeWidth
		} else {
			break
		}
	}

	content := string(text[:bytesCount])
	if isNumber {
		return NewToken(Number, content), bytesCount
	}
	return NewToken(Text, content), bytesCount
}

// isSpecialChar checks if a character is a special markdown character.
func isSpecialChar(c rune) bool {
	switch c {
	case '_', '*', '#', '`', '[', ']', '(', ')', '!', '?', '~', '-', '<', '>',
		'+', '.', '$', '=', '|', ':', '^', '\'', '\\', '/', '\n', ' ':
		return true
	default:
		return false
	}
}
