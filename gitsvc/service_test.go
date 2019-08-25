package gitsvc

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/config"
	"github.com/ujent/go-git-app/contract"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const userName = "test_user"
const userEmail = "test_user@gmail.com"
const remote = "http://35.239.165.218:9000/gitea/testrepo"
const remoteUser = "gitea@gitea.com"
const remotePsw = "secret123"

func TestPush(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	_, err = svc.CreateRemote(userName, r, remote, "")
	if err != nil {
		t.Fatal(err)
	}

	cr := &contract.Credentials{Name: remoteUser, Password: remotePsw}
	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	_, err = svc.Pull(rq, "", cr)
	if err != nil {
		t.Fatal(err)
	}

	fs, err := svc.Filesystem(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("Change hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit(rq, "add README")
	if err != nil {
		t.Fatal(err)
	}

	err = svc.Push(rq, "", cr)
	if err != nil {
		t.Error(err)
	}
}

func TestRepositories(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"
	r2 := "repo_2"

	err = svc.CreateRepository(userName, r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(userName, r1)

	err = svc.CreateRepository(userName, r2)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(userName, r2)

	repos, err := svc.Repositories(userName)
	if err != nil {
		t.Error(err)
	}

	must := 2
	if len(repos) < must {
		t.Errorf("Repositories quantity (%d) is less than %d \n", len(repos), must)
	}

	fmt.Println(repos)
}

func TestCreateRepository(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"

	err = svc.CreateRepository(userName, r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(userName, r1)

	repos, err := svc.Repositories(userName)
	if err != nil {
		t.Error(err)
	}

	must := 1
	if len(repos) < must {
		t.Fatalf("Repositories quantity (%d) is less than %d \n", len(repos), must)
	}

	if r1 != repos[0] {
		t.Errorf("Wrong repository name. Must: %s, has: %s\n", r1, repos[0])
	}
}

func TestRemoveRepository(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"

	err = svc.CreateRepository(userName, r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(userName, r1)

	err = svc.RemoveRepository(userName, r1)
	if err != nil {
		t.Error(err)
	}

	repos, err := svc.Repositories(userName)
	if err != nil {
		t.Error(err)
	}

	must := 0
	if len(repos) < must {
		t.Fatalf("Repositories quantity (%d) is less than %d \n", len(repos), must)
	}
}

func TestCurrentRepository(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"

	err = svc.CreateRepository(userName, r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(userName, r1)

	name := svc.CurrentRepository()

	if name != r1 {
		t.Errorf("Wrong current repository. Must: %s, has: %s\n", r1, name)
	}
}

func TestOpenRepository(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"

	err = svc.CreateRepository(userName, r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(userName, r1)

	r2 := "repo_2"

	err = svc.CreateRepository(userName, r2)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(userName, r2)

	err = svc.OpenRepository(userName, r1)
	if err != nil {
		t.Fatal(err)
	}

	name := svc.CurrentRepository()

	if name != r1 {
		t.Errorf("Wrong current repository. Must: %s, has: %s\n", r1, name)
	}
}

func TestCreateBranch(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit(rq, "add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CreateBranch(userName, r, b, "")
	if err != nil {
		t.Fatal(err)
	}

	branches, err := svc.Branches(userName, r)
	fmt.Println(branches)

	var hasBranch bool

	for _, br := range branches {
		if b == br {
			hasBranch = true
			break
		}
	}

	if !hasBranch {
		t.Errorf("Branch %s not found\n", b)
	}

	err = svc.RemoveRepository(userName, r)
	if err != nil {
		t.Error(err)
	}
}

//test checkout with existing branch
func TestCheckoutBranch1(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq1 := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq1, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	masterHash, err := svc.Commit(rq1, "add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CreateBranch(userName, r, b, "")
	if err != nil {
		t.Fatal(err)
	}

	f, err = fs.Create("Example.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("test"))

	rq2 := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: b}
	err = svc.Add(rq2, "Example.txt")

	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit(rq2, "add Example")
	if err != nil {
		t.Error(err)
	}

	mbr := "master"
	err = svc.CheckoutBranch(userName, r, mbr)
	if err != nil {
		t.Fatal(err)
	}

	current, err := svc.CurrentBranch()
	if err != nil {
		t.Fatal(err)
	}

	if current.Name != mbr {
		t.Errorf("Wrong current branch. Must: %s, has: %s \n", mbr, current.Name)
	}

	if current.Hash != masterHash {
		t.Errorf("Wrong branch hash. Must: %s, has: %s\n", masterHash, current.Hash)
	}
}

//test checkout with not existing branch
func TestCheckoutBranch2(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	brHash, err := svc.Commit(rq, "add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CheckoutBranch(userName, r, b)
	if err != nil {
		t.Fatal(err)
	}

	current, err := svc.CurrentBranch()
	if err != nil {
		t.Fatal(err)
	}

	if current.Name != b {
		t.Errorf("Wrong current branch. Must: %s, has: %s \n", b, current.Name)
	}

	if current.Hash != brHash {
		t.Errorf("Wrong branch hash. Must: %s, has: %s\n", brHash, current.Hash)
	}
}

func TestRemoveBranch(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit(rq, "add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CreateBranch(userName, r, b, "")
	if err != nil {
		t.Fatal(err)
	}

	branches, err := svc.Branches(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(branches)

	var hasBranch bool

	for _, br := range branches {
		if b == br {
			hasBranch = true
			break
		}
	}

	if !hasBranch {
		t.Errorf("Branch %s not found\n", b)
	}

	err = svc.RemoveBranch(userName, r, b)
	if err != nil {
		t.Fatal(err)
	}

	branches, err = svc.Branches(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(branches)

	hasBranch = false

	for _, br := range branches {
		if b == br {
			hasBranch = true
			break
		}
	}

	if hasBranch {
		t.Errorf("Branch %s wasn't deleted\n", b)
	}
}

func TestBranches(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit(rq, "add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CreateBranch(userName, r, b, "")
	if err != nil {
		t.Fatal(err)
	}

	branches, err := svc.Branches(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(branches)

	must := 2

	if len(branches) != must {
		t.Errorf("Wrong branches quantity. Must: %d, has: %d\n", must, len(branches))
	}
}

func TestCheckout(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	h1, err := svc.Commit(rq, "add README")
	if err != nil {
		t.Error(err)
	}

	f, err = fs.Create("Example.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("test"))

	err = svc.Add(rq, "Example.txt")
	if err != nil {
		t.Fatal(err)
	}

	rq.Branch = "master"
	_, err = svc.Commit(rq, "add Example")
	if err != nil {
		t.Error(err)
	}

	err = svc.Checkout(userName, r, h1)
	if err != nil {
		t.Fatal(err)
	}

	current, err := svc.CurrentBranch()
	if err != nil {
		t.Fatal(err)
	}

	if current.Hash != h1 {
		t.Errorf("Wrong branch hash. Must: %s, has: %s\n", h1, current.Hash)
	}
}

//test with master
func TestCurrentBranch1(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	h, err := svc.Commit(rq, "add README")
	if err != nil {
		t.Error(err)
	}

	current, err := svc.CurrentBranch()
	if err != nil {
		t.Fatal(err)
	}

	if current.Hash != h {
		t.Errorf("Wrong branch hash. Must: %s, has: %s\n", h, current.Hash)
	}

	must := "master"

	if current.Name != must {
		t.Errorf("Wrong branch name. Must: %s, has: %s\n", must, current.Name)

	}
}

//test with not master branch
func TestCurrentBranch2(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq1 := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq1, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	h, err := svc.Commit(rq1, "add README")
	if err != nil {
		t.Fatal(err)
	}

	br := "topic"

	err = svc.CreateBranch(userName, r, br, "")
	if err != nil {
		t.Fatal(err)
	}

	f, err = fs.Create("Example.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("test, go-git!"))

	rq2 := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: br}
	err = svc.Add(rq2, "Example.txt")

	if err != nil {
		t.Fatal(err)
	}

	h, err = svc.Commit(rq2, "add example")
	if err != nil {
		t.Fatal(err)
	}

	current, err := svc.CurrentBranch()
	if err != nil {
		t.Fatal(err)
	}

	if current.Hash != h {
		t.Errorf("Wrong branch hash. Must: %s, has: %s\n", h, current.Hash)
	}

	must := br

	if current.Name != must {
		t.Errorf("Wrong branch name. Must: %s, has: %s\n", must, current.Name)

	}
}

func TestLog(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	h1, err := svc.Commit(rq, "add README")
	if err != nil {
		t.Error(err)
	}

	f, err = fs.Create("Example.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("test"))

	rq.Branch = "master"
	err = svc.Add(rq, "Example.txt")
	if err != nil {
		t.Fatal(err)
	}

	h2, err := svc.Commit(rq, "add Example")
	if err != nil {
		t.Error(err)
	}

	log, err := svc.Log(rq)
	if err != nil {
		t.Fatal(err)
	}

	must := 2

	if len(log) != must {
		t.Errorf("Wrong log length. Must: %d, has: %d\n", must, len(log))
	}

	has1 := false
	has2 := false

	for _, c := range log {
		if c.Hash == h1 {
			has1 = true
			continue
		}

		if c.Hash == h2 {
			has2 = true
		}
	}

	if !has1 {
		t.Error("No first commit in logs")
	}

	if !has2 {
		t.Error("No second commit in logs")
	}
}

func TestCommit(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	rq := &contract.BaseRequest{User: &contract.User{Name: userName, Email: userEmail}, Repository: r, Branch: ""}
	fs, err := svc.Filesystem(userName, r)

	if err != nil {
		t.Fatal(err)
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add(rq, "README.md")
	if err != nil {
		t.Fatal(err)
	}

	msg := "add README"
	h, err := svc.Commit(rq, msg)
	if err != nil {
		t.Fatal(err)
	}

	rq.Branch = "master"

	log, err := svc.Log(rq)
	if err != nil {
		t.Fatal(err)
	}

	must := 1

	if len(log) != must {
		t.Errorf("Wrong log length. Must: %d, has: %d\n", must, len(log))
	}

	logCommit := log[0]

	if logCommit.Hash != h {
		t.Errorf("Wrong commit hash. Must: %s, has: %s\n", h, logCommit.Hash)
	}

	if logCommit.Message != msg {
		t.Errorf("Wrong commit message. Must: %s, has: %s\n", msg, logCommit.Message)
	}

	if logCommit.Author.Name != userName {
		t.Errorf("Wrong commit author name. Must: %s, has: %s\n", userName, logCommit.Author.Name)
	}

	if logCommit.Author.Email != userEmail {
		t.Errorf("Wrong commit author email. Must: %s, has: %s\n", userEmail, logCommit.Author.Email)
	}
}

func TestCreateRemote(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	_, err = svc.CreateRemote(userName, r, remote, "")
	if err != nil {
		t.Fatal(err)
	}

	rName := "origin"
	rem, err := svc.Remote(userName, r, rName)
	if err != nil {
		t.Fatal(err)
	}

	if len(rem.Config().URLs) != 1 {
		t.Fatalf("Wrong remote number. Must: 1, has: %d\n", len(rem.Config().URLs))
	}

	if rem.Config().URLs[0] != remote {
		t.Fatalf("Wrong remote url. Must: %s, has: %s\n", remote, rem.Config().URLs[0])
	}

	err = rem.Fetch(&git.FetchOptions{
		RemoteName: rName,
		Auth:       &http.BasicAuth{Username: remoteUser, Password: remotePsw},
	})

	if err != nil {
		t.Error(err)
	}
}

func TestRemoveRemote(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	_, err = svc.CreateRemote(userName, r, remote, "")
	if err != nil {
		t.Fatal(err)
	}

	rName := "origin"
	rem, err := svc.Remote(userName, r, rName)
	if err != nil {
		t.Fatal(err)
	}

	if len(rem.Config().URLs) != 1 {
		t.Fatalf("Wrong remote number. Must: 1, has: %d\n", len(rem.Config().URLs))
	}

	if rem.Config().URLs[0] != remote {
		t.Fatalf("Wrong remote url. Must: %s, has: %s\n", remote, rem.Config().URLs[0])
	}

	err = svc.RemoveRemote(userName, r, rName)
	if err != nil {
		t.Fatal(err)
	}

	rem, err = svc.Remote(userName, r, rName)
	if err != nil {
		if err != git.ErrRemoteNotFound {
			t.Error(err)
		}
	} else {
		t.Error("Remote wasn't removed")
	}
}

func TestRemote(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(userName, r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, r)

	_, err = svc.CreateRemote(userName, r, remote, "")
	if err != nil {
		t.Fatal(err)
	}

	rName := "origin"
	rem, err := svc.Remote(userName, r, rName)
	if err != nil {
		t.Fatal(err)
	}

	if len(rem.Config().URLs) != 1 {
		t.Fatalf("Wrong remote number. Must: 1, has: %d\n", len(rem.Config().URLs))
	}

	if rem.Config().URLs[0] != remote {
		t.Fatalf("Wrong remote url. Must: %s, has: %s\n", remote, rem.Config().URLs[0])
	}
}

func TestRemotes(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(s, db)
	if err != nil {
		t.Fatal(err)
	}

	repo := "repo_1"

	err = svc.CreateRepository(userName, repo)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(userName, repo)

	_, err = svc.CreateRemote(userName, repo, remote, "")
	if err != nil {
		t.Fatal(err)
	}

	rName := "origin"
	remotes, err := svc.Remotes(userName, repo)
	if err != nil {
		t.Fatal(err)
	}

	var rem *git.Remote

	for _, r := range remotes {
		if r.Config().Name == rName {
			rem = r
			break
		}
	}

	if rem == nil {
		t.Fatal("Remote wasn't found")
	}

	if len(rem.Config().URLs) != 1 {
		t.Fatalf("Wrong remote number. Must: 1, has: %d\n", len(rem.Config().URLs))
	}

	if rem.Config().URLs[0] != remote {
		t.Fatalf("Wrong remote url. Must: %s, has: %s\n", remote, rem.Config().URLs[0])
	}
}
