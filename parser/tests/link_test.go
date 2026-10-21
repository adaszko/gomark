package tests

import (
	"testing"

	"github.com/adaszko/gomark/ast"
	"github.com/adaszko/gomark/parser/internal"
	"github.com/adaszko/gomark/parser/tokenizer"
	"github.com/stretchr/testify/require"
)

func TestLinkParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "[](https://example.com)",
			node: &ast.Link{
				Content: []ast.Node{},
				URL:     "https://example.com",
			},
		},
		{
			text: "! [](https://example.com)",
			node: nil,
		},
		{
			text: "[your/slash](https://example.com)",
			node: &ast.Link{
				Content: []ast.Node{&ast.Text{Content: "your/slash"}},
				URL:     "https://example.com",
			},
		},
		{
			text: "[hello world](https://example.com)",
			node: &ast.Link{
				Content: []ast.Node{&ast.Text{Content: "hello world"}},
				URL:     "https://example.com",
			},
		},
		{
			text: "[hello world](https://example.com)",
			node: &ast.Link{
				Content: []ast.Node{&ast.Text{Content: "hello world"}},
				URL:     "https://example.com",
			},
		},
		{
			text: `[\[link\]](https://example.com)`,
			node: &ast.Link{
				Content: []ast.Node{&ast.EscapingCharacter{Symbol: "["}, &ast.Text{Content: `link`}, &ast.EscapingCharacter{Symbol: "]"}},
				URL:     "https://example.com",
			},
		},
		{
			text: "[![alt text](https://example.com/image.png)](https://example.com)",
			node: &ast.Link{
				Content: []ast.Node{&ast.Image{AltText: "alt text", URL: "https://example.com/image.png"}},
				URL:     "https://example.com",
			},
		},
		{
			text: "[![](https://example.com/image.png)](https://example.com)",
			node: &ast.Link{
				Content: []ast.Node{&ast.Image{AltText: "", URL: "https://example.com/image.png"}},
				URL:     "https://example.com",
			},
		},
	}
	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := internal.NewLinkParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}

func TestLinkSpaces(t *testing.T) {
	text := "[](foo://bar baz)"
	expected_node := &ast.Link{
		Content: []ast.Node{},
		URL:     "foo://bar baz",
	}
	tokens := tokenizer.Tokenize(text)
	actual_node, _ := internal.NewLinkParser().Match(tokens)
	require.Equal(t, expected_node, actual_node)
}
