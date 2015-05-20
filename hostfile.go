package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/howeyc/gopass"
)

//HostFileParser used for parsing CSV files
type HostFileParser struct {
	Hosts []HostProfile
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
								hProfile.Password = hp.PasswordPrompt(hProfile.Host)
							} else {
								hProfile.Password = line[l]
							}
						case 3:
							hProfile.Key = line[l]
						}
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
	Username string
	Password string
	Host     string
	Key      string
}
