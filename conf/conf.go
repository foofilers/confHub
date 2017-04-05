package conf

import (
	cnf "github.com/spf13/viper"
	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"fmt"
	"strings"
)

var (
	initialized = false
)

func IsInitialized() bool{
	return initialized
}

func InitConfFromFile(configName string, configPaths ...string) {
	if initialized {
		log.Warn("configuration already initialized, overwriting it")
	}
	fmt.Println("Initializing configuration")
	cnf.SetConfigName(configName)
	cnf.SetConfigType("yaml")
	for _, cnfPath := range (configPaths) {
		cnf.AddConfigPath(cnfPath)
	}
	err := cnf.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	cnf.WatchConfig()
	cnf.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Configuration changed")
	})
	initialized = true
}

func InitConf(yamlContent string) {
	if initialized {
		log.Warn("configuration already initialized, overwriting it")
	}
	cnf.SetConfigType("yaml")
	if err := cnf.ReadConfig(strings.NewReader(yamlContent)); err != nil {
		log.Fatal(err)
	}
	initialized = true
}
