package main

import (
	"fmt"

	"github.com/Juniper/go-netconf/netconf"
)

//Agent agent to connect and issue commands to hosts
type Agent struct {
	Username   string
	Password   string
	Key        string
	Host       string
	Session    *netconf.Session
	MsgChannel chan Message
	parser     Parser
}

//Dial connect to host
func (a *Agent) Dial() {
	var err error
	if a.Username != "" && a.Password != "" {
		a.Session, err = netconf.DialSSH(a.Host, netconf.SSHConfigPassword(a.Username, a.Password))
		if err != nil {
			a.returnMsg("", "", err)
		}
	}
}

//Close close session to host
func (a *Agent) Close() {
	a.Session.Close()
}

func (a *Agent) returnMsg(data string, command string, err error) {
	a.MsgChannel <- Message{Host: a.Host, Command: command, Data: data, Error: err}
}

//RunCommand Run a command against a host
func (a *Agent) RunCommand(command string) {
	reply, err := a.Session.Exec(netconf.RawMethod(fmt.Sprintf("<command format=\"ascii\">%s</command>", command)))
	if err != nil {
		a.returnMsg(reply.Data, command, err)
	}
	v := a.parser.Trim(reply.Data)
	a.returnMsg(v, command, nil)
}
