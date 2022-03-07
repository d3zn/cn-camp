package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Service struct {
		Port int    `json:"port" yaml:"port"`
		Mode string `json:"mode" yaml:"mode"`
	} `json:"service" yaml:"service"`
	Log struct {
		Level   string `json:"level" yaml:"level"`
		LogFile string `json:"logfile" yaml:"logfile"`
	} `json:"log" yaml:"log"`
}

func LoadConfig() (*Config, error) {
	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(curPath)
	if err != nil {
		return nil, err
	}

	confInstance := new(Config)
	confFile := curPath + "/conf/conf.yaml"

	bs, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, fmt.Errorf("read config file failed, err: %v", err)
	}

	err = yaml.Unmarshal(bs, confInstance)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed, err:%v", err)
	}
	if err != nil {
		return nil, err
	}
	return confInstance, nil
}
