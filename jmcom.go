package main

import (
	"flag"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

var msgChannel chan Message
var ctrlChannel chan Message
var ctrlChans map[string]chan Message
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
	flag.Parse()
	//Check variables
	//Spawn agent per connection
	msgChannel = make(chan Message)
	ctrlChannel = make(chan Message)
	ctrlChans = make(map[string]chan Message)

	go func() {
		for {
			select {
			case msg, chanOpen := <-msgChannel:
				if chanOpen && msg.Error != nil {
					//handle error
					log.Errorf("Session %d error: %s", msg.SessionID, msg.Error)
				} else if chanOpen && msg.Data != "" && msg.Host != "" {
					log.Printf("Host %s SessionID %d\n%s", msg.Host, msg.SessionID, msg.Data)
					wg.Done()
				} else {
					log.Println("CLOSED")
					return
				}
			}
		}
	}()

	if hosts != "" && user != "" && password != "" {
		hs := strings.Split(hosts, ",")
		for _, v := range hs {
			ctrlChans[v] = make(chan Message)
			a := &Agent{Username: user, Password: password, Host: v, CtrlChannel: ctrlChans[v], MsgChannel: msgChannel}
			wg.Add(1)
			log.Println("Connecting to", v)
			go a.Run()
		}
	}
	time.Sleep(10 * time.Second)
	//Run command against hosts
	log.Println(command)
	for item := range ctrlChans {
		ctrlChans[item] <- Message{Command: command}
	}

	//return results
	wg.Wait()
	close(msgChannel)
	close(ctrlChannel)
	for item := range ctrlChans {
		close(ctrlChans[item])
	}
	log.Println("Complete")
}
