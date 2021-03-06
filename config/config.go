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
	GRPCDialAddress string
	//TBotToken - telegram bot token
	TBotToken string
	//SuperAdmin - telegram username of superadmin
	SuperAdmin []string
	//PermFIle - permissions filename
	PermFile string
)

const (
	modeKey         = "MODE"
	devMode         = "development"
	grpcDialAddrKey = "DIALADDR"
	telegramToken   = "TTOKEN"
	superAdmin      = "SADMIN"
	permFile        = "PERMFILE"
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
	GRPCDialAddress = os.Getenv(grpcDialAddrKey)
	if GRPCDialAddress == "" {
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
	PermFile = os.Getenv(permFile)
	if PermFile == "" {
		log.Fatal("permission filename must be set")
	}
}
