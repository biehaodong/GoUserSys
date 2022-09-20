package utils

import (
	"fmt"
	"time"

	"gopkg.in/ini.v1"
)

var (
	AppMode        string
	ServerPort     string
	HTTPServerPort string
	ClientPoolSize int

	Db              string
	DbHost          string
	DbPort          string
	DbUser          string
	DbPassWord      string
	DbName          string
	DbAddress       string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime = 14000 * time.Second

	RedisAddress  string
	RedisPoolSize int
	TokenLife     int
	KeyLife       int

	StaticFilePath string
	DefaultImage   string

	TCPServerLogPath  string
	HTTPServerLogPath string
)

func init() {
	file, err := ini.Load("/Users/haodong.bie/GolandProjects/GoUserManaSys/config/config.ini")
	if err != nil {
		fmt.Println("配置文件错误:", err)
	}
	loadServer(file)
	loadData(file)
	loadRedis(file)
	loadStatic(file)
	loadLog(file)
}
func loadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	ServerPort = file.Section("server").Key("HttpPort").MustString(":3000")
	ClientPoolSize, _ = file.Section("server").Key("ClientPoolSize").Int()
	HTTPServerPort = file.Section("server").Key("HTTPServerPort").MustString("1806")

}

func loadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("dao")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("root")
	DbName = file.Section("database").Key("DbName").MustString("userInfo")
	DbAddress = file.Section("database").Key("DbAddress").MustString(
		"root:root@tcp_server.log(127.0.0.1:3306)/userInfo")
	MaxIdleConns, _ = file.Section("database").Key("MaxIdleConns").Int()
	MaxOpenConns, _ = file.Section("database").Key("MaxOpenConns").Int()
}
func loadRedis(file *ini.File) {
	RedisAddress = file.Section("redis").Key("RedisAddress").MustString("localhost:6379")
	RedisPoolSize, _ = file.Section("redis").Key("RedisPoolSize").Int()
	TokenLife, _ = file.Section("redis").Key("TokenLife").Int()
	KeyLife, _ = file.Section("redis").Key("KeyLife").Int()
}
func loadStatic(file *ini.File) {
	StaticFilePath = file.Section("static").Key("StaticFilePath").MustString("./static/")
	DefaultImage = file.Section("static").Key("DefaultImage").MustString("girl.jpg")
}
func loadLog(file *ini.File) {
	TCPServerLogPath = file.Section("log").Key("hTCPServerLogPath").MustString("./log/tcp_server.log")
	HTTPServerLogPath = file.Section("log").Key("HTTPServerLogPath").MustString("./log/http_server.log")
}
