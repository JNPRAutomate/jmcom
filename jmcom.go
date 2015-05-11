package main

import (
	"flag"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var ctrlChannel map[string]chan Message
var hosts string
var user string
var password string
var command string
var wg sync.WaitGroup

func init() {

	//TODO: Add CSV and SSH key support

	//Define flags for calling script
	flag.StringVar(&hosts, "hosts", "", "Define hosts to connect to: 1.2.3.3 or 2.3.4.5,1.2.3.4")
	flag.StringVar(&user, "user", "", "Specify the username to use against hosts")
	flag.StringVar(&password, "password", "", "Specify password to use with hosts")
	flag.StringVar(&command, "command", "", "Command to run against hosts")
	//offer logging to single or multiple files
}

func main() {
	//Check variables
	//Spawn agent per connection

	if hosts != "" && user != "" && password != "" {
		hs := strings.Split(hosts, ",")
		for _, v := range hs {
			ctrlChannel[v] = make(chan Message)
			a := &Agent{Username: user, Password: password, Host: v, CtrlChannel: ctrlChannel[v]}
			wg.Add(1)
			log.Println("Connecting to ", v)
			go a.Run()
		}
	}
	//Run command against hosts

	//return results
}
