package contract

import "time"

//CredentialsPayload base user information for operations which requires auth
type CredentialsPayload struct {
	Name string `json:"name"`
	Psw  string `json:"psw"`
}

//RepoRQ is the request payload for operations with repository
type RepoRQ struct {
	Name string `json:"name"`
}

//RepositoriesRS - the response to repositories request
type RepositoriesRS struct {
	Repos []RepoRS `json:"repos"`
}

//RepoRS - base info about repository
type RepoRS struct {
	Name      string `json:"name"`
	IsCurrent bool   `json:"isCurrent"`
}

//BranchRQ is the request payload for operations with repository
type BranchRQ struct {
	Name string `json:"name"`
}

//BranchesRS - the response to branches request
type BranchesRS struct {
	Branches []BranchRS `json:"branches"`
}

//BranchRS - base info about branch
type BranchRS struct {
	Name      string `json:"name"`
	IsCurrent bool   `json:"isCurrent"`
}

//UserRS -  common user information
type UserRS struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	IsCurrent bool   `json:"isCurrent"`
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

//FileInfoRS - common information about files in repository
type FileInfoRS struct {
	Path       string `json:"path"`
	IsConflict bool   `json:"isConflict"`
}

//FilesRS - the response to files request
type FilesRS struct {
	Files []FileInfoRS
}

//CloneRQ is the request payload for clone repository
type CloneRQ struct {
	Auth     *CredentialsPayload `json:"auth,omitempty"`
	URL      string              `json:"URL"`
	RepoName string              `json:"repoName"`
}

// CommitRQ - request for commit operation
type CommitRQ struct {
	Message string `json:"message"`
}

// PullRQ - request for pull operation
type PullRQ struct {
	Auth   *CredentialsPayload `json:"auth,omitempty"`
	Remote string              `json:"remote"`
}

// PushRQ - request for push operation
type PushRQ struct {
	Auth   *CredentialsPayload `json:"auth,omitempty"`
	Remote string              `json:"remote"`
}

// MergeRQ - request for merge operation
type MergeRQ struct {
	Branch string `json:"branch"`
}

//MsgResult - common result returns message
type MsgResult struct {
	Msg string `json:"msg"`
}

//SwitchUserRQ - request for changing user from which we are using app
type SwitchUserRQ struct {
	Name string `json:"name"`
}

//UsersRS - response for users request
type UsersRS struct {
	Users []UserRS `json:"users"`
}
