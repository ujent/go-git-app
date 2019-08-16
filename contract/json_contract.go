package contract

//CredentialsPayload base user information for operations which requires auth
type CredentialsPayload struct {
	Name string
	Psw  string
}

//RepoRQ is the request payload for operations with repository
type RepoRQ struct {
	Name string `json:"name"`
}

//RepositoriesRS - tge response to repositories request
type RepositoriesRS struct {
	Names []string `json:"names"`
}
