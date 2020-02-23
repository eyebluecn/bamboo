package core

const (
	//authentication key of cookie
	COOKIE_AUTH_KEY = "_ak"

	USERNAME_KEY = "_username"
	PASSWORD_KEY = "_password"

	DEFAULT_SERVER_PORT = 6020

	//db table's prefix. bamboo10_ means current version is bamboo:1.0.x
	TABLE_PREFIX = "bamboo10_"

	VERSION = "1.0.0"
)

type Config interface {
	Installed() bool
	ServerPort() int
	//get the mysql url. eg. bamboo:bamboo123@tcp(127.0.0.1:3306)/bamboo?charset=utf8&parseTime=True&loc=Local
	MysqlUrl() string
	//when installed by user. Write configs to bamboo.json
	FinishInstall(mysqlPort int, mysqlHost string, mysqlSchema string, mysqlUsername string, mysqlPassword string)
}
