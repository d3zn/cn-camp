package service

import (
	"fmt"
	"os"
	"time"
)

var _ Service = &service{}

type Service interface {
	Header()
	Env(key string) (string, error)
	Health() error
}

type service struct {
	startTime time.Time
}

func (s *service) Header() {
	time.Sleep(60 * time.Second)
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

func (s *service) Health() error {
	d := time.Since(s.startTime)
	if d.Seconds() > 60 {
		return fmt.Errorf("server down!")
	}
	return nil
}

func NewService(startTime time.Time) *service {
	return &service{
		startTime: startTime,
	}
}
