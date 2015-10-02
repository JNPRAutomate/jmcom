package main

// Agent an interface for agents
type Agent interface {
	// Run start the agent to run commands
	Run()
	// RunCommand run a specific command
	RunCommand(string)
	// Close close the connection to the agent
	Close()
	// Dial connect to the agent
	Dial()
	// dial
	dial()
	returnMsg(string, string, error)
}
