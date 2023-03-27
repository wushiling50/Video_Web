package conf

import (
	"fmt"
	"main/Video_Web/model"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	Host      string
	VideoPath string
	ImgPath   string

	ValidEmail string
	SmtpHost   string
	SmtpEmail  string
	SmtpPass   string
)

func Init() {
	file, err := ini.Load("./Video_Web/conf/conf.ini")
	if err != nil {
		fmt.Println("配置文件读取有误,请检查配置文件路径")
	}
	LoadServer(file)
	LoadMysql(file)
	LoadPath(file)
	LoadEmail(file)

	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=True"}, "")
	model.Database(path)

}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("DB").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()

}

func LoadPath(file *ini.File) {
	Host = file.Section("path").Key("Host").String()
	VideoPath = file.Section("path").Key("VideoPath").String()
	ImgPath = file.Section("path").Key("ImgPath").String()
}

func LoadEmail(file *ini.File) {
	ValidEmail = file.Section("email").Key("ValidEmail").String()
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmtpEmail = file.Section("email").Key("SmtpEmail").String()
	SmtpPass = file.Section("email").Key("SmtpPass").String()
}
