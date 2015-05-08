package main

import "strings"

type Parser struct {
}

func (p *Parser) Trim(input string) string {
	t := strings.Trim(strings.Trim(strings.TrimSpace(input), "<output>"), "</output>")
	return t
}
