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
	Name  string
	Email string
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
