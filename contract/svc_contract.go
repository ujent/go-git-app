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
	GitRoot    string
	FsType     FsType
}

// BaseRequest - rq for most git operations
type BaseRequest struct {
	User       *User
	Repository string
	Branch     string
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

//FileInfo - common information about files in repository
type FileInfo struct {
	Path       string
	IsConflict bool
}

//FsType - concrete type of filesystem
type FsType int

const (
	FsTypeInvalid FsType = 0
	FsTypeMySQL   FsType = 1
	FsTypeLocal   FsType = 2
)

func ToFsType(t int) FsType {

	switch t {
	case 1:
		return FsTypeMySQL
	case 2:
		return FsTypeLocal
	default:
		return FsTypeInvalid
	}
}
