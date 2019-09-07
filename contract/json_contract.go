package contract

import "time"

//CredentialsPayload base user information for operations which requires auth
type CredentialsPayload struct {
	Name string `json:"name"`
	Psw  string `json:"psw"`
}

//RepoRQ is the request payload for operations with repository
type RepoRQ struct {
	User string `json:"user"`
	Repo string `json:"repo"`
}

//RepositoriesRQ - the request for repositories operation
type RepositoriesRQ struct {
	User string `json:"user"`
}

//RepositoriesRS - the response to repositories request
type RepositoriesRS struct {
	Repos []string `json:"repos"`
}

//RepoRS - base info about repository
type RepoRS struct {
	Name string `json:"name"`
}

//BranchRQ is the request payload for operations with repository
type BranchRQ struct {
	Branch string `json:"branch"`
	User   string `json:"user"`
	Repo   string `json:"repo"`
}

//BranchesRQ - the request for branches operation
type BranchesRQ struct {
	User       string `json:"user"`
	Repository string `json:"repo"`
}

//BranchesRS - the response to branches request
type BranchesRS struct {
	Branches []string `json:"branches"`
	Current  string   `json:"current"`
}

//BranchRS - base info about branch
type BranchRS struct {
	Name string `json:"name"`
}

//UserRS -  common user information
type UserRS struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

//LogRS - the response to log request
type LogRS struct {
	Commits []CommitRS `json:"commits"`
}

//CommitRS - base commit information
type CommitRS struct {
	Author  *UserRS   `json:"author"`
	Hash    string    `json:"hash"`
	Message string    `json:"msg"`
	Date    time.Time `json:"date"`
}

type FileRQ struct {
	Base       *BaseRequestRQ `json:"base"`
	Path       string         `json:"path"`
	IsConflict bool           `json:"isConflict"`
}

type FileRS struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type AddFileRQ struct {
	Base    *BaseRequestRQ `json:"base"`
	Path    string         `json:"path"`
	Content string         `json:"content"`
}

type RemoveFileRQ struct {
	Base *BaseRequestRQ `json:"base"`
	Path string         `json:"path"`
}

type EditFileRQ struct {
	Base    *BaseRequestRQ `json:"base"`
	Path    string         `json:"path"`
	Content string         `json:"content"`
}

//FileInfoRS - common information about files in repository
type FileInfoRS struct {
	Path       string `json:"path"`
	IsConflict bool   `json:"isConflict"`
}

//FilesRQ - the files request
type FilesRQ struct {
	Base *BaseRequestRQ `json:"base"`
}

//FilesRS - the response to files request
type FilesRS struct {
	Files []FileInfoRS `json:"files"`
}

//CloneRQ is the request payload for clone repository
type CloneRQ struct {
	User string              `json:"user"`
	Auth *CredentialsPayload `json:"auth,omitempty"`
	URL  string              `json:"URL"`
}

// CommitRQ - request for commit operation
type CommitRQ struct {
	Base    *BaseRequestRQ `json:"base"`
	Message string         `json:"message"`
}

// PullRQ - request for pull operation
type PullRQ struct {
	Base   *BaseRequestRQ      `json:"base"`
	Auth   *CredentialsPayload `json:"auth,omitempty"`
	Remote string              `json:"remote"`
}

// PushRQ - request for push operation
type PushRQ struct {
	Base   *BaseRequestRQ      `json:"base"`
	Auth   *CredentialsPayload `json:"auth,omitempty"`
	Remote string              `json:"remote"`
}

// MergeRQ - request for merge operation
type MergeRQ struct {
	Base   *BaseRequestRQ `json:"base"`
	Theirs string         `json:"theirs"`
}

type MergeRS struct {
	Message       string `json:"msg"`
	IsFastforward bool   `json:"isFF"`
}

type AbortMergeRQ struct {
	Base *BaseRequestRQ `json:"base"`
}

//MsgResult - common result returns message
type MsgResult struct {
	Msg string `json:"msg"`
}

//SwitchUserRQ - request for changing user from which we are using app
type SwitchUserRQ struct {
	Name string `json:"name"`
}

// BaseRequestRQ - rq for most git operations
type BaseRequestRQ struct {
	User       string `json:"user"`
	Repository string `json:"repo"`
	Branch     string `json:"branch"`
}
