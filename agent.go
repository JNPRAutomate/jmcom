package main

import (
	"errors"
	"fmt"

	"github.com/Juniper/go-netconf/netconf"
)

//Agent agent to connect and issue commands to hosts
type Agent struct {
	SessionID   int
	HostProfile *HostProfile
	Session     *netconf.Session
	CtrlChannel chan Message
	MsgChannel  chan Message
	parser      Parser
}

//Run set agent to run commands
func (a *Agent) Run() {
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

//Dial connect to host
func (a *Agent) dial() error {
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

//Close close session to host
func (a *Agent) Close() {
	a.Session.Close()
}

func (a *Agent) returnMsg(data string, command string, err error) {
	a.MsgChannel <- Message{Host: a.HostProfile.Host, SessionID: a.SessionID, Command: command, Data: data, Error: err}
}

//RunCommand Run a command against a host
func (a *Agent) RunCommand(command string) {
	reply, err := a.Session.Exec(netconf.RawMethod(fmt.Sprintf("<command format=\"ascii\">%s</command>", command)))
	if err != nil {
		a.returnMsg(reply.Data, command, err)
	}
	v := a.parser.Trim(reply.Data)
	a.returnMsg(v, command, err)
}
