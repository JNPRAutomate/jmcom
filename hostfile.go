package main

//HostFileParser used for parsing CSV files
type HostFileParser struct {
	Hosts []HostProfile
}

//PasswordPrompt prompt a user for an interactive password
func (hp *HostFileParser) PasswordPrompt() {

}

//HostProfile used as a profile for host configurations
type HostProfile struct {
	Username string
	Password string
	Host     string
	Key      string
}
