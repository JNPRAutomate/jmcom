package main

import (
	"errors"
	"fmt"

	"github.com/Juniper/go-netconf/netconf"
)

const (
	// CfgTypeNull the default type for CfgType
	CfgTypeNull = 0
)

// CfgType configuration type
type CfgType int

const (
	// JunosAgentModeNull the default mode for an JunosAgent
	JunosAgentModeNull = 0
	// JunosAgentModeOp enters the JunosAgent into operational mode
	JunosAgentModeOp = 1
	// JunosAgentModeConfig enters the JunosAgent into configuration mode
	JunosAgentModeConfig = 2
)

// JunosAgentMode the mode in which the JunosAgent is running
type JunosAgentMode int

// JunosAgent JunosAgent to connect and issue commands to hosts
type JunosAgent struct {
	// SessionID the ID receieved from the netconf session
	// This is the PID that you are running as on the device
	SessionID int
	// HostProfile the connection information for authentication to the host
	HostProfile *HostProfile
	// Session the current netconf session for the host
	Session *netconf.Session
	// CtrlChannel channel that is used to control the JunosAgent
	CtrlChannel chan Message
	// MsgChannel channel that is used to send messages
	MsgChannel chan Message
	// parser the parser for dealing with response data
	parser Parser
	// Mode the mode the JunosAgent is operating in
	Mode JunosAgentMode
}

// Run set JunosAgent to run commands
func (a *JunosAgent) Run() {
	err := a.dial()
	if err != nil {
		a.returnMsg("", "", err)
		return
	}
	a.returnMsg("", "", nil)
	for {
		select {
		case msg, chanOpen := <-a.CtrlChannel:
			if chanOpen {
				a.RunCommand(msg.Command)
			} else {
				a.Close()
				return
			}
		}
	}
}

// Dial connect to the remote host
func (a *JunosAgent) Dial() error {
	return a.dial()
}

// dial connect to host
func (a *JunosAgent) dial() error {
	var err error
	if a.HostProfile.Username != "" && len(a.HostProfile.GetSSHClientConfig().Auth) > 0 {
		a.Session, err = netconf.DialSSH(a.HostProfile.Host, a.HostProfile.GetSSHClientConfig())
		if err != nil {
			return err
		}
		a.SessionID = a.Session.SessionID
		return nil
	}
	err = errors.New("Host Profile incorrectly defined")
	return err
}

// Close close session to host
func (a *JunosAgent) Close() {
	a.Session.Close()
}

// returnMsg send back a return message to the host
func (a *JunosAgent) returnMsg(data string, command string, err error) {
	a.MsgChannel <- Message{Host: a.HostProfile.Host, SessionID: a.SessionID, Command: command, Data: data, Error: err}
}

// RunCommand Run a command against a host
func (a *JunosAgent) RunCommand(command string) {
	r, err := a.Session.Exec(netconf.RawMethod(fmt.Sprintf("<command format=\"ascii\">%s</command>", command)))
	if err != nil {
		a.returnMsg(r.Data, command, err)
	}
	v := a.parser.Trim(r.Data)
	a.returnMsg(v, command, err)
}

// LoadConfig load a set of configuration data to the device
func (a *JunosAgent) LoadConfig(config string, cfgType CfgType, overwrite bool) {
	// pull the config to the agent
	// open the configuration mode
	// commit the configuration
	r, err := a.Session.Exec(netconf.RawMethod(fmt.Sprintf("", config)))
	if err != nil {
		a.returnMsg(r.Data, "config", err)
	}
	a.returnMsg(r.Data, "config", err)
}
