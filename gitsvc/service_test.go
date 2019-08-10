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

	err = svc.CreateRepository(r2)
	if err != nil {
		t.Error(err)
	}

	repos, err := svc.Repositories()
	if err != nil {
		t.Error(err)
	}

	must := 2
	if len(repos) < must {
		t.Errorf("Repositories quantity (%d) is less than %d \n", len(repos), must)
	}

	fmt.Println(repos)

	err = svc.RemoveRepository(r1)
	if err != nil {
		t.Error(err)
	}

	err = svc.RemoveRepository(r2)
	if err != nil {
		t.Error(err)
	}
}
