package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/howeyc/gopass"
)

// HostFileParser used for parsing CSV files
type HostFileParser struct {
	GlobalPassword string
	GlobalKey      string
	Hosts          []HostProfile
}

// Parse parse the hostfile
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
				// host,username,password,keyfile
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
				// load the global key and or password if not methids are specified
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

// PasswordPrompt prompt a user for an interactive password
func (hp *HostFileParser) PasswordPrompt(hostname string) string {
	fmt.Printf("Enter password for %s: ", hostname)
	text := gopass.GetPasswdMasked()
	return string(text)
}
