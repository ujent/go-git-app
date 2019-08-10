package contract

import "time"

//ServerSettings - common server settings
type ServerSettings struct {
	Port       string
	GitConnStr string
	GitRemote  string
}

//User -  current user information
type User struct {
	Name  string
	Email string
}

//Credentials - user credentials
type Credentials struct {
	Name     string
	Password string
}

//Commit - base commit information
type Commit struct {
	Author  *User
	Hash    string
	Message string
	Date    time.Time
}
