package rest

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

func Start(addr string, wg *sync.WaitGroup, quitCh chan bool) {
	logrus.Info("starting Rest")
	wg.Add(1)
	go func() {
		select {
		case <-quitCh:
			Stop()
			wg.Done()
		}
	}()
	if viper.GetString("environment") == "prod" || viper.GetString("environment") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	if viper.GetString("environment") == "test" {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()
	InitApplication(server)
	go func() {
		server.Run(addr)
	}()
}
func Stop() {
	logrus.Info("Stopping rest")
}