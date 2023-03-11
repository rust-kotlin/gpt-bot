package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Apikey string `json:"apikey"`
	Model  string `json:"model"`
	Proxy  string `json:"proxy"`
}

// LoadConfig 从文件中读取配置，如果文件不存在则返回默认配置
func LoadConfig(filename string) (*Config, error) {
	config := &Config{
		Apikey: "your api key",
		Model:  "gpt-3.5-turbo",
		Proxy:  "http://localhost:7890",
	}

	// 如果文件不存在，则创建文件并写入默认配置
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := SaveConfig(filename, config); err != nil {
			return nil, err
		}
		panic("Please edit the config.json file and restart the program.")
	}

	// 读取文件并解析为配置项
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(file, config); err != nil {
		return nil, err
	}
	return config, nil
}

// SaveConfig 将配置项写入文件
func SaveConfig(filename string, config *Config) error {
	file, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, file, 0644); err != nil {
		return err
	}
	return nil
}
