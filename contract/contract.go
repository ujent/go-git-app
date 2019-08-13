package contract

import (
	"errors"
	"io"
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

// FileStage during merge
type FileStage int

const (
	//Unexpected - is unexpected index.Stage
	Unexpected FileStage = -1
	// Merged is the default stage, fully merged
	Merged FileStage = 0
	// AncestorMode is the base revision
	AncestorMode FileStage = 1
	// OurMode is the first tree revision, ours
	OurMode FileStage = 2
	// TheirMode is the second tree revision, theirs
	TheirMode FileStage = 3
)

//MergeFile - file during merge
type MergeFile struct {
	Path   string
	Stage  FileStage
	Reader io.Reader
}
