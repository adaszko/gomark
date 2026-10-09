package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		text   string
		tokens []*Token
	}{
		{
			text: "*Hello world!",
			tokens: []*Token{
				{
					Type:  Asterisk,
					Value: "*",
				},
				{
					Type:  Text,
					Value: "Hello",
					Pos:   1,
				},
				{
					Type:  Space,
					Value: " ",
					Pos:   6,
				},
				{
					Type:  Text,
					Value: "world",
					Pos:   7,
				},
				{
					Type:  ExclamationMark,
					Value: "!",
					Pos:   12,
				},
			},
		},
		{
			text: `# hello 
 world`,
			tokens: []*Token{
				{
					Type:  PoundSign,
					Value: "#",
				},
				{
					Type:  Space,
					Value: " ",
					Pos:   1,
				},
				{
					Type:  Text,
					Value: "hello",
					Pos:   2,
				},
				{
					Type:  Space,
					Value: " ",
					Pos:   7,
				},
				{
					Type:  NewLine,
					Value: "\n",
					Pos:   8,
				},
				{
					Type:  Space,
					Value: " ",
					Pos:   9,
				},
				{
					Type:  Text,
					Value: "world",
					Pos:   10,
				},
			},
		},
	}

	for _, test := range tests {
		result := Tokenize(test.text)
		require.Equal(t, test.tokens, result)
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		tokens []*Token
		sep    TokenType
		result [][]*Token
	}{
		{
			tokens: []*Token{
				{
					Type:  Asterisk,
					Value: "*",
				},
				{
					Type:  Text,
					Value: "Hello",
				},
				{
					Type:  Space,
					Value: " ",
				},
				{
					Type:  Text,
					Value: "world",
				},
				{
					Type:  ExclamationMark,
					Value: "!",
				},
			},
			sep: Asterisk,
			result: [][]*Token{
				{},
				{
					{
						Type:  Text,
						Value: "Hello",
					},
					{
						Type:  Space,
						Value: " ",
					},
					{
						Type:  Text,
						Value: "world",
					},
					{
						Type:  ExclamationMark,
						Value: "!",
					},
				},
			},
		},
	}

	for _, test := range tests {
		result := Split(test.tokens, test.sep)
		for index, tokens := range result {
			require.Equal(t, Stringify(test.result[index]), Stringify(tokens))
		}
	}
}

func TestGetFirstLine(t *testing.T) {
	tests := []struct {
		tokens []*Token
		want   []*Token
	}{
		{
			tokens: []*Token{
				{
					Type:  Asterisk,
					Value: "hello world",
				},
				{
					Type:  NewLine,
					Value: "\n",
				},
			},
			want: []*Token{
				{
					Type:  Asterisk,
					Value: "hello world",
				},
			},
		},
	}

	for _, test := range tests {
		result := GetFirstLine(test.tokens)
		require.Equal(t, test.want, result)
	}
}
