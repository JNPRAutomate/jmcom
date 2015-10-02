package main

// Agent an interface for agents
type Agent interface {
	// Run start the agent to run commands
	Run()
	// RunCommand run a specific command
	RunCommand(string)
	// Close close the connection to the agent
	Close()
	// Dial connect to the remote host
	Dial() error
	// dial connection or where all the good stuff happens
	dial() error
	// returnMsg  return a response back to the message channel
	returnMsg(string, string, error)
}
