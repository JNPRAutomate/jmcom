package main

import (
	"fmt"
	"sync"

	"github.com/Juniper/go-netconf/netconf"
	log "github.com/Sirupsen/logrus"
)

//Agent agent to connect and issue commands to hosts
type Agent struct {
	SessionID   int
	Username    string
	Password    string
	Key         string
	Host        string
	Session     *netconf.Session
	CtrlChannel chan Message
	MsgChannel  chan Message
	parser      Parser
	connectWg   sync.WaitGroup
}

//Run set agent to run commands
func (a *Agent) Run() {
	a.Dial()
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
func (a *Agent) Dial() {
	var err error
	if a.Username != "" && a.Password != "" {
		a.Session, err = netconf.DialSSH(a.Host, netconf.SSHConfigPassword(a.Username, a.Password))
		if err != nil {
			a.returnMsg("", "", err)
			return
		}
		a.SessionID = a.Session.SessionID
		log.Infoln("Connected to", a.Host)
		connectWg.Done()
	} else if a.Username != "" && a.Password == "" && a.Key != "" {
		//setup key based auth
	}
}

//Close close session to host
func (a *Agent) Close() {
	a.Session.Close()
}

func (a *Agent) returnMsg(data string, command string, err error) {
	a.MsgChannel <- Message{Host: a.Host, SessionID: a.SessionID, Command: command, Data: data, Error: err}
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
