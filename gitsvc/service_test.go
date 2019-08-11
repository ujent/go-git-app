package gitsvc

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/config"
	"github.com/ujent/go-git-app/contract"
)

const userName = "test_user"
const userEmail = "test_user@gmail.com"

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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"
	r2 := "repo_2"

	err = svc.CreateRepository(r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(r1)

	err = svc.CreateRepository(r2)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(r2)

	repos, err := svc.Repositories()
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"

	err = svc.CreateRepository(r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(r1)

	repos, err := svc.Repositories()
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"

	err = svc.CreateRepository(r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(r1)

	err = svc.RemoveRepository(r1)
	if err != nil {
		t.Error(err)
	}

	repos, err := svc.Repositories()
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"

	err = svc.CreateRepository(r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(r1)

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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r1 := "repo_1"

	err = svc.CreateRepository(r1)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(r1)

	r2 := "repo_2"

	err = svc.CreateRepository(r2)
	if err != nil {
		t.Error(err)
	}

	defer svc.RemoveRepository(r2)

	err = svc.OpenRepository(r1)
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(r)

	fs := svc.Filesystem()
	if fs == nil {
		t.Fatal("No filesystem")
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit("add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CreateBranch(b, "")
	if err != nil {
		t.Fatal(err)
	}

	branches, err := svc.Branches()
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

	err = svc.RemoveRepository(r)
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(r)

	fs := svc.Filesystem()
	if fs == nil {
		t.Fatal("No filesystem")
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}

	masterHash, err := svc.Commit("add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CreateBranch(b, "")
	if err != nil {
		t.Fatal(err)
	}

	f, err = fs.Create("Example.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("test"))

	err = svc.Add("Example.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit("add Example")
	if err != nil {
		t.Error(err)
	}

	mbr := "master"
	err = svc.CheckoutBranch(mbr)
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(r)

	fs := svc.Filesystem()
	if fs == nil {
		t.Fatal("No filesystem")
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}

	brHash, err := svc.Commit("add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CheckoutBranch(b)
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(r)

	fs := svc.Filesystem()
	if fs == nil {
		t.Fatal("No filesystem")
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit("add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CreateBranch(b, "")
	if err != nil {
		t.Fatal(err)
	}

	branches, err := svc.Branches()
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

	err = svc.RemoveBranch(b)
	if err != nil {
		t.Fatal(err)
	}

	branches, err = svc.Branches()
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(r)

	fs := svc.Filesystem()
	if fs == nil {
		t.Fatal("No filesystem")
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit("add README")
	if err != nil {
		t.Error(err)
	}

	b := "test_branch"

	err = svc.CreateBranch(b, "")
	if err != nil {
		t.Fatal(err)
	}

	branches, err := svc.Branches()
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

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(r)

	fs := svc.Filesystem()
	if fs == nil {
		t.Fatal("No filesystem")
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}

	h1, err := svc.Commit("add README")
	if err != nil {
		t.Error(err)
	}

	f, err = fs.Create("Example.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("test"))

	err = svc.Add("Example.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Commit("add Example")
	if err != nil {
		t.Error(err)
	}

	err = svc.Checkout(h1)
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

func TestCurrentBranch(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.GitConnStr)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	svc, err := New(&contract.User{Name: userName, Email: userEmail}, s, db)
	if err != nil {
		t.Fatal(err)
	}

	r := "repo_1"

	err = svc.CreateRepository(r)
	if err != nil {
		t.Fatal(err)
	}

	defer svc.RemoveRepository(r)

	fs := svc.Filesystem()
	if fs == nil {
		t.Fatal("No filesystem")
	}

	f, err := fs.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("hello, go-git!"))

	err = svc.Add("README.md")
	if err != nil {
		t.Fatal(err)
	}

	h, err := svc.Commit("add README")
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
}
