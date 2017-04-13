package ws

import (
	"gopkg.in/kataras/iris.v6/adaptors/websocket"
	"sync"
	"github.com/Sirupsen/logrus"
)

var watchConnMutex sync.Mutex

var watchConn = make(map[string]map[string]websocket.Connection) //applications,connID

func AddWatchConnAndStartNotifier(appName string, conn websocket.Connection) {
	logrus.Debugf("connection %s watch  application %s", conn.ID(), appName)
	watchConnMutex.Lock()
	defer watchConnMutex.Unlock()
	var currConns map[string]websocket.Connection
	var found bool
	currConns, found = watchConn[appName]
	if !found || len(currConns) == 0 {
		currConns = make(map[string]websocket.Connection)
		watchConn[appName] = currConns
		//new app start notifier in background
		go StartNotifier(appName)
	}
	currConns[conn.ID()] = conn
}

func RemoveWatchAppConn(appName string, connID string) {
	watchConnMutex.Lock()
	defer watchConnMutex.Unlock()
	currConns, found := watchConn[appName]
	if !found {
		return
	}
	delete(currConns, connID)
	if len(currConns) == 0 {
		StopNotifier(appName)
	}
}

func RemoveWatchConn(connID string) {
	for appName, conns := range watchConn {
		for watchConnId := range conns {
			if watchConnId == connID {
				RemoveWatchAppConn(appName, connID)
			}
		}
	}
}

func GetWatchAppConns(appName string) map[string]websocket.Connection {
	return watchConn[appName]
}