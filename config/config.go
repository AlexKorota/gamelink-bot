package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

var (
	//DiaAddress - on following address client should dial to server
	DialAddress string
	//TBotToken - telegram bot token
	TBotToken string
	//SuperAdmin - telegram username of superadmin
	SuperAdmin []string
	//mongoAddrKey - network addres for mongoDB
	MongoAddr string
	//mongoDBName
	MongoDBName string
	//LogFile - name of log file
	LogFile string
)

const (
	modeKey       = "MODE"
	devMode       = "development"
	dialAddrKey   = "DIALADDR"
	telegramToken = "TTOKEN"
	superAdmin    = "SADMIN"
	mongoAddr     = "MONGOADDR"
	mongoDBName   = "MONGODBNAME"
	logFile       = "LOGFILE"
)

func init() {
	LoadEnvironment()
}

//GetEnvironment - this function returns mode string of the os environment or "development" mode if empty or not defined
func GetEnvironment() string {
	var env string
	if env = os.Getenv(modeKey); env == "" {
		return devMode
	}
	return env
}

//IsDevelopmentEnv - this function try to get mode environment and check it is development
func IsDevelopmentEnv() bool { return GetEnvironment() == devMode }

func LoadEnvironment() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = godotenv.Load(path.Join(wd, strings.ToLower(GetEnvironment())+".env"))
	if err != nil {
		log.Warning(err.Error())
	}
	DialAddress = os.Getenv(dialAddrKey)
	if DialAddress == "" {
		log.Fatal("server address must be set")
	}
	TBotToken = os.Getenv(telegramToken)
	if TBotToken == "" {
		log.Fatal("telegram token must be set")
	}
	SA := os.Getenv(superAdmin)
	if SA == "" {
		log.Fatal("should be at least one super admin")
	}
	SuperAdmin = strings.Split(SA, ",")

	MongoAddr = os.Getenv(mongoAddr)
	if MongoAddr == "" {
		log.Fatal("mongo address must be set")
	}
	MongoDBName = os.Getenv(mongoDBName)
	if MongoAddr == "" {
		log.Fatal("mongo DB name must be set")
	}
	LogFile = os.Getenv(logFile)
	if LogFile == "" {
		log.Fatal("log filename must be set")
	}
}
