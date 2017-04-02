package log

import (
	"os"
	log "github.com/Sirupsen/logrus"
)

func InitLog(){
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Debug("Log initialized")
}