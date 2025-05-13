package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Port      int    `yaml:"port"`
	ServerUrl string `yaml:"server_url"`
}

var GlobalConfig *Config

func ParseConfig(filename string) *Config {
	// 先尝试从当前工作目录加载
	currentDir, err := os.Getwd()
	if err == nil {
		configPath := filepath.Join(currentDir, filename)
		data, err := os.ReadFile(configPath)
		if err == nil {
			var cfg *Config
			err = yaml.Unmarshal(data, &cfg)
			if err == nil {
				GlobalConfig = cfg
				log.Printf("Config loaded from: %s\n", configPath)
				return cfg
			}
		}
	}

	// 如果从当前目录加载失败，再尝试从可执行文件目录加载
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("\n\nUnable to get executable path: %+v\n\n", err)
	}

	dir := filepath.Dir(executable)
	configPath := filepath.Join(dir, filename)

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("\n\nNot found config file, %+v\n\n", err)
	}

	var cfg *Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("\n\nUnable to parse config file, %+v\n\n", err)
	}
	GlobalConfig = cfg
	log.Printf("Config loaded from: %s\n", configPath)
	return cfg
}
