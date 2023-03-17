package jsonconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Apikey        string  `json:"apikey"`
	Model         string  `json:"model"`
	Proxy         string  `json:"proxy"`
	MaxTokens     int     `json:"max_tokens"`
	SuperUsers    []int64 `json:"super_users"`
	WeatherApikey string  `json:"weather_apikey"`
	Temperature   float32 `json:"temperature"`
	BaseUrl       string  `json:"base_url"`
}

// LoadConfig 从文件中读取配置，如果文件不存在则返回默认配置
func LoadConfig(filename string) (*Config, error) {
	config := &Config{
		Apikey:        "your api key",
		Model:         "gpt-3.5-turbo",
		Proxy:         "http://localhost:7890",
		MaxTokens:     800,
		SuperUsers:    []int64{},
		WeatherApikey: "",
		Temperature:   0.5,
		BaseUrl:       "",
	}

	// 如果文件不存在，则创建文件并写入默认配置
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := SaveConfig(filename, config); err != nil {
			return nil, err
		}
		return nil, err
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
