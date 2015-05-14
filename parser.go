package main

import "strings"

//Parser used to parse and massage the output from commands
type Parser struct {
}

//Trim remove unwanted values from the output
func (p *Parser) Trim(input string) string {
	t := strings.Replace(input, "<output>", "", -1)
	t1 := strings.Replace(t, "</output>", "", -1)
	return t1
}
