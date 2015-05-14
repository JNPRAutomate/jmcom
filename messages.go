package main

type Message struct {
	Host      string
	SessionID int
	Data      string
	Command   string
	Error     error
}
