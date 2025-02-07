package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	BaseURL  string `json:"baseURL"`
	DeviceID string `json:"deviceID"`
	Token    string `json:"token"`
}

func main() {
	// 检查本地是否有账户状态
	// 创建文件
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// 读取文件
	buffer, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	// 解析文件
	var account Account
	if err := json.Unmarshal(buffer, &account); err != nil {
		log.Fatal(err)
	}
}
