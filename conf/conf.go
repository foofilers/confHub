package conf

import (
	cnf "github.com/spf13/viper"
	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"fmt"
)

func InitConf() {
	fmt.Println("Initializing configuration")
	cnf.SetConfigName("confHub")
	cnf.SetConfigType("yaml")
	cnf.AddConfigPath("./conf")
	err := cnf.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	cnf.WatchConfig()
	cnf.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Configuration changed")
	})
}
