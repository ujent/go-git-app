package contract

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
