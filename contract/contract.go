package contract

import (
	"errors"
	"time"
)

//ErrGitRepositoryNotSet - occurs when repository wasn't chosen
var ErrGitRepositoryNotSet = errors.New("Git repository isn't set")

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

//Branch - base information about branch
type Branch struct {
	Name string
	Hash string
}
