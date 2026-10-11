package internal

import (
	"github.com/adaszko/gomark/ast"
	"github.com/adaszko/gomark/parser/tokenizer"
)

type LineBreakParser struct{}

func NewLineBreakParser() *LineBreakParser {
	return &LineBreakParser{}
}

func (*LineBreakParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	if len(tokens) == 0 {
		return nil, 0
	}
	if tokens[0].Type != tokenizer.NewLine {
		return nil, 0
	}
	return &ast.LineBreak{}, 1
}
