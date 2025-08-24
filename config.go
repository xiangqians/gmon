// @author xiangqian
// @date 2025/07/27 17:54
package main

import (
	"gmon/pkg/prom"
	pkg_ini "gopkg.in/ini.v1"
	"strings"
)

// LoadConfig 加载配置文件
func LoadConfig() (Config, error) {
	file, err := pkg_ini.Load("config.ini")
	if err != nil {
		return Config{}, err
	}

	// http
	section, err := file.GetSection("http")
	if err != nil {
		return Config{}, err
	}
	var http = Http{
		Port:   uint16(section.Key("port").MustUint()),
		Prefix: strings.TrimSpace(section.Key("prefix").String()),
		User:   strings.TrimSpace(section.Key("user").String()),
		Passwd: strings.TrimSpace(section.Key("passwd").String()),
	}

	// prom
	section, err = file.GetSection("prom")
	if err != nil {
		return Config{}, err
	}
	var prom = prom.Config{
		Host: strings.TrimSpace(section.Key("host").String()),
		Port: uint16(section.Key("port").MustUint()),
	}

	return Config{Http: http, Prom: prom}, nil
}

// Config 配置
type Config struct {
	Http Http        // HTTP 配置
	Prom prom.Config // Prometheus 配置
}

// Http HTTP 配置
type Http struct {
	Port   uint16 // 监听端口
	Prefix string // HTTP 请求前缀
	User   string // 登录用户
	Passwd string // 登录密码（如果含有特殊字符，如 #，则使用反引号括起来）
}
