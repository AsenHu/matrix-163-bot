package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Path    string
	Content Content
}

type Content struct {
	Matrix  Matrix  `json:"matrix"`
	Netease Netease `json:"netease"`
}

type Matrix struct {
	BaseURL  string `json:"baseURL"`
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceID string `json:"deviceID"`
	Token    string `json:"token"`
}

type Netease struct {
	MUSIC_U  string `json:"MUSIC_U"`
	Quailty  string `json:"quailty"`
	CheckMd5 bool   `json:"checkMd5"`
}

func NewConfig(path string) Config {
	return Config{
		Path: path,
		Content: Content{
			Matrix: Matrix{
				BaseURL:  "https://example.com",
				Username: "@bot:example.com",
				Password: "password",
			},
			Netease: Netease{
				Quailty:  "higher",
				CheckMd5: true,
			},
		},
	}
}

func (c *Config) Load() (err error) {
	// 读取文件
	buffer, err := os.ReadFile(c.Path)
	if err != nil {
		return
	}
	// 解析文件
	err = json.Unmarshal(buffer, &c.Content)
	if err != nil {
		return
	}
	return
}

func (c *Config) Save() (err error) {
	// 序列化文件
	buffer, err := json.MarshalIndent(c.Content, "", "  ")
	if err != nil {
		return
	}
	// 写入文件
	err = os.WriteFile(c.Path, buffer, 0600)
	if err != nil {
		return
	}
	return
}
