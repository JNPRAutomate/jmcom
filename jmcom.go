package main

import (
	"flag"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var msgChannel chan Message
var ctrlChans map[string]chan Message
var hosts string
var user string
var password string
var commands string
var commandFile string
var logs bool
var logLocation string
var hostsFile string
var commandWg sync.WaitGroup
var connectWg sync.WaitGroup
var recWg sync.WaitGroup

func init() {
	//TODO: Add CSV and SSH key support

	//Define flags for calling script
	flag.StringVar(&hosts, "hosts", "", "Define hosts to connect to: 1.2.3.3 or 2.3.4.5,1.2.3.4")
	flag.StringVar(&user, "user", "", "Specify the username to use against hosts")
	flag.StringVar(&password, "password", "", "Specify password to use with hosts")
	flag.StringVar(&commands, "command", "", "Commands to run against host: \"show version\" or for multiple commands \"show version\",\"show chassis hardware\"")
	flag.StringVar(&hostsFile, "hosts-file", "", "File to load hosts from")
	flag.StringVar(&commandFile, "cmd-file", "", "File to load commands from")
	flag.BoolVar(&logs, "log", false, "Log output for each host to a seperate file")
	flag.StringVar(&logLocation, "logdir", "", "Directory to write logs to. Default is current directory")
	//offer logging to single or multiple files
}

func main() {
	flag.Parse()
	//Check variables
	//Spawn agent per connection
	msgChannel = make(chan Message)
	ctrlChans = make(map[string]chan Message)

	recWg.Add(1)
	go func() {
		for {
			select {
			case msg, chanOpen := <-msgChannel:
				if chanOpen && msg.Error != nil {
					log.Errorf("Session %d error: %s", msg.SessionID, msg.Error)
				} else if chanOpen && msg.Data != "" && msg.Host != "" {
					log.Printf("Host %s SessionID %d\n%s", msg.Host, msg.SessionID, msg.Data)
					commandWg.Done()
				} else {
					recWg.Done()
					return
				}
			}
		}
	}()

	if hosts != "" && user != "" && password != "" {
		hs := strings.Split(hosts, ",")
		for _, v := range hs {
			ctrlChans[v] = make(chan Message)
			connectWg.Add(1)
			a := &Agent{Username: user, Password: password, Host: v, connectWg: connectWg, CtrlChannel: ctrlChans[v], MsgChannel: msgChannel}
			log.Println("Connecting to", v)
			go a.Run()
		}
	}
	connectWg.Wait()
	//Run command against hosts
	cmds := strings.Split(commands, ",")
	for _, v := range cmds {
		for item := range ctrlChans {
			commandWg.Add(1)
			ctrlChans[item] <- Message{Command: v}
		}
	}

	//return results
	commandWg.Wait()
	close(msgChannel)
	for item := range ctrlChans {
		close(ctrlChans[item])
	}
	recWg.Wait()
	log.Println("Complete")
}
