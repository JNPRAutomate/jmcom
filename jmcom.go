package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/howeyc/gopass"
)

//Tool globals
var msgChannel chan Message
var ctrlChans map[string]chan Message
var logFiles map[string]*os.File
var cmds []string
var hostps []*HostProfile

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

	//prompt for password if not defined
	if passPrompt && password == "" {
		password = promptPassword()
	}

	//ME list
	// hosts sshKey password hostsFile
	if hosts != "" && hostsFile != "" {
		log.Infoln("Combining command line hosts with host from host file")
	}

	if sshKey != "" && password != "" {
		log.Infoln("Using both ssh key and password as auth methods")
	}

	//MB list
	//commands and commandFile
	if commands != "" && commandFile != "" {
		log.Infoln("Combining command line commands with command file command set")
	}

	//create channels for communication
	msgChannel = make(chan Message)
	ctrlChans = make(map[string]chan Message)
	//Create map for logging
	logFiles = make(map[string]*os.File)

	//Split hosts

	//Host file parsing
	hfp := &HostFileParser{GlobalPassword: password, GlobalKey: sshKey}

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
		cmds = append(cmds, strings.Split(string(cmdFile), "\n")...)
	}

	if commands != "" {
		cmds = append(cmds, strings.Split(commands, ",")...)
	}

	//setup hosts from file
	if hostsFile != "" {
		h, err := hfp.Parse(hostsFile)
		if err != nil {
			log.Fatalf("Unable to parse host file: %s", err)
		}
		hostps = append(hostps, h...)
		for i := range h {
			fmt.Printf("%#v\n", h[i])
		}
		hostps = append(hostps, h...)
	}

	//setup hosts
	if hosts != "" {
		clihosts := strings.Split(hosts, ",")
		for i := range clihosts {
			hp := &HostProfile{Host: clihosts[i], Username: user}
			if password != "" {
				hp.LoadPassword(password)
			}

			if sshKey != "" {
				err := hp.LoadKey(sshKey)
				if err != nil {
					log.Fatalf("Unable to load key specified by flag: %s", err)
				}
			}
			hostps = append(hostps, hp)
		}
	}

	//setup log files
	if logs {
		for _, v := range hostps {
			var err error
			logFiles[v.Host], err = OpenLog(logLocation, v.Host)
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
					log.Errorf("Host: %s error: %s", msg.Host, msg.Error)
					if msg.SessionID == 0 && msg.Command == "" {
						connectWg.Done()
					}
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

	if len(hostps) > 0 {
		for i := range hostps {
			ctrlChans[hostps[i].Host] = make(chan Message)
			connectWg.Add(1)
			a := &Agent{HostProfile: hostps[i], connectWg: connectWg, CtrlChannel: ctrlChans[hostps[i].Host], MsgChannel: msgChannel}
			log.Println("Connecting to", hostps[i].Host)
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
	if len(ctrlChans) == 0 {
		//no connections initiated, the user needs help
		flag.PrintDefaults()
		return
	}
	log.Println("Tasks Complete")
}

func promptPassword() string {
	fmt.Printf("Enter password: ")
	text := gopass.GetPasswdMasked()
	return string(text)
}
