package main

import "strings"

//Parser used to parse and massage the output from commands
type Parser struct {
}

//Trim remove unwanted values from the output
func (p *Parser) Trim(input string) string {
	t := strings.Replace(strings.Replace(input, "<output>", "", -1), "</output>", "", -1)
	t1 := strings.Replace(strings.Replace(t, "<configuration-information>", "", -1), "</configuration-information>", "", -1)
	return strings.Replace(strings.Replace(t1, "<configuration-output>", "", -1), "</configuration-output>", "", -1)
}
