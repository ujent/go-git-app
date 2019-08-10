package gitSvc

import (
	"fmt"
	"testing"

	"github.com/ujent/go-git-app/config"
	"github.com/ujent/go-git-app/contract"
	"github.com/ujent/go-git-app/gitsvc"
)

const userName = "test_user"
const userEmail = "test_user@gmail.com"

func TestRepositories(t *testing.T) {
	s, err := config.ParseTest()
	if err != nil {
		t.Fatal(err)
	}

	svc, err := gitsvc.New(&contract.Credentials{Name: userName, Email: userEmail}, s)
	if err != nil {
		t.Fatal(err)
	}

	err = svc.CreateRepository("repo_1")
	if err != nil {
		t.Error(err)
	}

	err = svc.CreateRepository("repo_2")
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
}
