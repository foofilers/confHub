package main

import (
	"github.com/foofilers/cfhd/rpc"
	"github.com/foofilers/cfhd/db"
	"github.com/foofilers/cfhd/conf"
	"github.com/sirupsen/logrus"
	"github.com/foofilers/cfhd/rest"
	"sync"
	"os"
	"os/signal"
	"syscall"
)

var quitGrpcCh chan bool
var quitRestCh chan bool

var wg sync.WaitGroup

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("cfhd starting")
	conf.Init()
	db.Init()
	quitGrpcCh = make(chan bool)
	quitRestCh = make(chan bool)
	initKiller()
	rpc.Start("0.0.0.0:50051", &wg, quitGrpcCh)
	rest.Start("0.0.0.0:8080", &wg,quitRestCh)
	wg.Wait()
	logrus.Info("cfhd terminated")

}

func initKiller() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == syscall.SIGINT {
				logrus.Info("Received kill signal")
				quitGrpcCh <- true
				quitRestCh <- true
				close(quitGrpcCh)
				close(quitRestCh)
			}
		}
	}()
}

