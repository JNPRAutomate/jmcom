package main

import (
	"io/ioutil"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// HostProfile used as a profile for host configurations
type HostProfile struct {
	Username    string
	Host        string
	AuthMethods []ssh.AuthMethod
}

// LoadKey load SSH key into auth methods
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

// LoadPassword load a password into auth methods
func (p *HostProfile) LoadPassword(password string) {
	p.AuthMethods = append(p.AuthMethods, ssh.Password(password))
}

// GetSSHClientConfig generates an ssh.ClientConfig to be use for SSH authentication
func (p *HostProfile) GetSSHClientConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{User: p.Username, Auth: p.AuthMethods}
}
