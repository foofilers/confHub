package conf

import (
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

func Init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("/etc/cfhd/")   // path to look for the config file in
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		logrus.Errorf("Error loading configuration:%s",err)
	}
}