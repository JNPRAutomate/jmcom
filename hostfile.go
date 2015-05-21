package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/howeyc/gopass"
)

//HostFileParser used for parsing CSV files
type HostFileParser struct {
	GlobalPassword string
	GlobalKey      string
	Hosts          []HostProfile
}

//Parse parse
func (hp *HostFileParser) Parse(file string) ([]*HostProfile, error) {
	hostP := []*HostProfile{}
	fl, err := filepath.Abs(file)
	if err != nil {
		return hostP, err
	}
	hostFile, err := ioutil.ReadFile(fl)
	if err != nil {
		return hostP, err
	}
	hostFiles := strings.Split(string(hostFile), "\n")
	for i := range hostFiles {
		if len(hostFiles[i]) > 1 {
			if string(hostFiles[i][0]) != "#" {
				line := strings.Split(hostFiles[i], ",")
				//host,username,password,keyfile
				hProfile := &HostProfile{}
				for l := range line {
					if line[l] != "" {
						switch l {
						case 0:
							hProfile.Host = line[l]
						case 1:
							hProfile.Username = line[l]
						case 2:
							if line[l] == "!!PROMPT!!" {
								hProfile.LoadPassword(hp.PasswordPrompt(hProfile.Host))
							} else {
								hProfile.LoadPassword(line[l])
							}
						case 3:
							hProfile.LoadKey(line[l])
						}
					}
				}
				if len(hProfile.AuthMethods) == 0 {
					if hp.GlobalKey != "" {
						hProfile.LoadKey(hp.GlobalKey)
					}
					if hp.GlobalPassword != "" {
						hProfile.LoadPassword(hp.GlobalPassword)
					}
				}
				hostP = append(hostP, hProfile)
			}
		}
	}
	return hostP, nil
}

//PasswordPrompt prompt a user for an interactive password
func (hp *HostFileParser) PasswordPrompt(hostname string) string {
	fmt.Printf("Enter password for %s: ", hostname)
	text := gopass.GetPasswdMasked()
	return string(text)
}

//HostProfile used as a profile for host configurations
type HostProfile struct {
	Username    string
	Host        string
	AuthMethods []ssh.AuthMethod
}

//LoadKey load SSH key into auth methods
func (p *HostProfile) LoadKey(key string) error {
	fl, err := filepath.Abs(key)
	if err != nil {
		return err
	}
	f, err := ioutil.ReadFile(fl)
	if err != nil {
		return err
	}
	k, err := ssh.ParsePrivateKey(f)
	if err != nil {
		return err
	}
	p.AuthMethods = append(p.AuthMethods, ssh.PublicKeys(k))
	return nil
}

//LoadPassword load a password into auth methods
func (p *HostProfile) LoadPassword(password string) {
	p.AuthMethods = append(p.AuthMethods, ssh.Password(password))
}

//GetSSHClientConfig generates an ssh.ClientConfig to be use for SSH authentication
func (p *HostProfile) GetSSHClientConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{User: p.Username, Auth: p.AuthMethods}
}
