package contract

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
