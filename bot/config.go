package bot

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ConfigDatabase struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Token    string         `yaml:"token"`
	Database ConfigDatabase `yaml:"database"`
}

func (ctx *Context) LoadConfig() error {
	file, err := os.OpenFile(ctx.ConfigFile.Filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(file)
	_ = file.Close()
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, &ctx.ConfigFile.Data)
}

func (ctx *Context) SaveConfig() error {
	bs, err := yaml.Marshal(ctx.ConfigFile.Data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ctx.ConfigFile.Filename, bs, 0644)
}
