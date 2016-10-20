package main

import (
	"fmt"

	"github.com/Juniper/go-netconf/netconf"
)

func foo() {
	username := "xxxxx"
	password := "xxxxxx"
	host := "172.19.100.49"

	s, err := netconf.DialSSH(host, netconf.SSHConfigPassword(username, password))
	if err != nil {
		panic(err)
	}

	defer s.Close()

	fmt.Printf("Session Id: %d\n\n", s.SessionID)

	reply, err := s.Exec(netconf.RawMethod("<command format=\"ascii\">show version</command>"))
	if err != nil {
		panic(err)
	}
	p := &Parser{}
	v := p.Trim(reply.Data)
	fmt.Println("REPLY START")
	fmt.Printf("%s", v)
	fmt.Println("REPLY END")
}
