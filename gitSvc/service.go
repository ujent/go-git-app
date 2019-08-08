package gitsvc

//Service - provides go-git functionality
type Service interface {
	CreateRepository(name string) error
	RemoveRepository(name string) error
	Branchs() error
	Fetch() error
	Pull() error
	Push() error
	Commit() error
	Merge(branch string) error
	Checkout(branch string) error
	Add() error
}

type service struct {
}

func NewSvc() Service {
	return &service{}
}

func (svc *service) CreateRepository(name string) error {
	return nil
}

func (svc *service) RemoveRepository(name string) error {
	return nil
}

func (svc *service) Branchs() error {
	return nil
}

func (svc *service) Fetch() error {
	return nil
}

func (svc *service) Pull() error {
	return nil
}

func (svc *service) Push() error {
	return nil
}

func (svc *service) Commit() error {
	return nil
}

func (svc *service) Merge(branch string) error {
	return nil
}

func (svc *service) Checkout(branch string) error {
	return nil
}

func (svc *service) Add() error {
	return nil
}
