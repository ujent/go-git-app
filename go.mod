module github.com/ujent/go-git-app

go 1.12

require (
	github.com/gliderlabs/ssh v0.2.2 // indirect
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20190725054713-01f96b0aa0cd // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/ujent/go-git-mysql v0.0.0-20190807142715-fbc9b2c784de

	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa // indirect
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.2 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
)

replace gopkg.in/src-d/go-git.v4 v4.13.1 => github.com/ujent/go-git v0.0.0-20190801043737-fd24d52a153b
