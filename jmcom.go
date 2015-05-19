package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

//Tool globals
var msgChannel chan Message
var ctrlChans map[string]chan Message
var logFiles map[string]*os.File
var cmds []string

//flags
var hosts string
var user string
var password string
var sshKey string
var commands string
var commandFile string
var logs bool
var passPrompt bool
var hostsFile string
var logLocation string

//Wait Groups for syncing go routines
var commandWg sync.WaitGroup
var connectWg sync.WaitGroup
var recWg sync.WaitGroup

func init() {
	//TODO: Add CSV and SSH key support

	//Define flags for calling script
	flag.StringVar(&hosts, "hosts", "", "Define hosts to connect to: 1.2.3.3 or 2.3.4.5,1.2.3.4")
	flag.StringVar(&user, "user", "", "Specify the username to use against hosts")
	flag.StringVar(&password, "password", "", "Specify password to use with hosts")
	flag.StringVar(&sshKey, "key", "", "Specify SSH key to use")
	flag.BoolVar(&passPrompt, "prompt", false, "Promots the user to enter a password interactively")
	flag.StringVar(&commands, "command", "", "Commands to run against host: \"show version\" or for multiple commands \"show version\",\"show chassis hardware\"")
	flag.StringVar(&hostsFile, "hosts-file", "", "File to load hosts from")
	flag.StringVar(&commandFile, "cmd-file", "", "File to load commands from")
	flag.BoolVar(&logs, "log", false, "Log output for each host to a seperate file")
	flag.StringVar(&logLocation, "logdir", "", "Directory to write logs to. Default is current directory")
}

func main() {
	flag.Parse()

	//create channels for communication
	msgChannel = make(chan Message)
	ctrlChans = make(map[string]chan Message)
	//Create map for logging
	logFiles = make(map[string]*os.File)

	//Split hosts
	hs := strings.Split(hosts, ",")
	cmds = strings.Split(commands, ",")

	//Host file parsing
	hp := &HostFileParser{}
	h := []*HostProfile{}

	//setup hosts from file
	if hostsFile != "" {
		var err error
		h, err = hp.Parse(hostsFile)
		if err != nil {
			log.Fatalf("Unable to parse host file: %s", err)
		}
	}

	//setup command file
	if commandFile != "" {
		fl, err := filepath.Abs(commandFile)
		if err != nil {
			log.Fatalln(err)
		}
		cmdFile, err := ioutil.ReadFile(fl)
		if err != nil {
			log.Fatalln(err)
		}
		cmds = strings.Split(string(cmdFile), "\n")
	}

	//setup log files
	if logs {
		for _, v := range hs {
			var err error
			logFiles[v], err = OpenLog(logLocation, v)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	recWg.Add(1)
	go func() {
		for {
			select {
			case msg, chanOpen := <-msgChannel:
				if chanOpen && msg.Error != nil {
					log.Errorf("Session %d error: %s", msg.SessionID, msg.Error)
				} else if chanOpen && msg.Data != "" && msg.Host != "" {
					if logs {
						log.SetOutput(logFiles[msg.Host])
						log.Printf("Host: %s SessionID: %d Command: %s\n%s", msg.Host, msg.SessionID, msg.Command, msg.Data)
						log.SetOutput(os.Stdout)
					} else {
						log.Printf("Host: %s SessionID: %d Command: %s\n%s", msg.Host, msg.SessionID, msg.Command, msg.Data)
					}
					commandWg.Done()
				} else {
					recWg.Done()
					return
				}
			}
		}
	}()

	if hosts != "" && user != "" && password != "" {
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
	for _, v := range cmds {
		c := strings.Replace(v, "\n", "", -1)
		if len(c) > 3 {
			for item := range ctrlChans {
				commandWg.Add(1)
				ctrlChans[item] <- Message{Command: c}
			}
		}
	}

	//return results
	commandWg.Wait()
	close(msgChannel)
	for item := range ctrlChans {
		close(ctrlChans[item])
	}
	recWg.Wait()
	log.Println("Tasks Complete")
}
