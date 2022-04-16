package service

import (
	"fmt"
	"m02/metrics"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var _ Service = &service{}

type Service interface {
	Header()
	Env(key string) (string, error)
	Health() error
	Random() (string, error)
}

type service struct {
	startTime time.Time
	crash     bool
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
	if !s.crash {
		return nil
	}
	d := time.Since(s.startTime)
	if d.Seconds() > 60 {
		return fmt.Errorf("server down!")
	}
	return nil
}

func (s *service) Random() (string, error) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	delay := rand.Intn(2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	return strconv.Itoa(delay) + "ms", nil
}

func NewService(startTime time.Time, crash bool) *service {
	return &service{
		startTime: startTime,
		crash:     crash,
	}
}
