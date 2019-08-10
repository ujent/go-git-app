package contract

import "time"

//ServerSettings - common server settings
type ServerSettings struct {
	Port       string
	GitConnStr string
	GitRemote  string
}

//Credentials -  current user information
type Credentials struct {
	Name  string
	Email string
}

//Commit - base commit information
type Commit struct {
	Author  *Credentials
	Hash    string
	Message string
	Date    time.Time
}
