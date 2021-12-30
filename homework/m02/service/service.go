package service

import (
	"fmt"
	"os"
)

var _ Service = &service{}

type Service interface {
	Header()
	Env(key string) (string, error)
	Health() string
}

type service struct{}

func (s *service) Header() {
	return
}

func (s *service) Env(key string) (string, error) {
	if err := os.Setenv("myEnv", "abcd"); err != nil {
		return "", err
	}
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("not found, key: %s ", key)
	}
	return v, nil
}

func (s *service) Health() string {
	return "health"
}

func NewService() *service {
	return &service{}
}
