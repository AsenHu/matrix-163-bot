package config

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Path    string
	Content Account
}

type Account struct {
	BaseURL  string `json:"baseURL"`
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceID string `json:"deviceID"`
	Token    string `json:"token"`
}

func NewConfig(path string) Config {
	return Config{
		Path: path,
	}
}

func (c *Config) Load() (err error) {
	// 读取文件
	buffer, err := os.ReadFile(c.Path)
	if (err != nil) {
		return
	}
	// 解析文件
	err = json.Unmarshal(buffer, &c.Content)
	if (err != nil) {
		return
	}
	return
}

func (c *Config) Save() (err error) {
	// 序列化文件
	buffer, err := json.MarshalIndent(c.Content, "", "  ")
	if (err != nil) {
		return
	}
	// 写入文件
	err = os.WriteFile(c.Path, buffer, 0600)
	if (err != nil) {
		return
	}
	return
}

func (c Config) Print() {
	log.Info().Str("BaseURL", c.Content.BaseURL).Msg("Account Info")
	log.Info().Str("Username", c.Content.Username).Msg("Account Info")
	log.Info().Str("Password", c.Content.Password).Msg("Account Info")
	log.Info().Str("DeviceID", c.Content.DeviceID).Msg("Account Info")
	log.Info().Str("Token", c.Content.Token).Msg("Account Info")
}
