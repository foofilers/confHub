package ws

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/etcd"
	cnf "github.com/spf13/viper"
	"context"
	"github.com/coreos/etcd/clientv3"
	"time"
)

var notifiers = make(map[string]*etcd.EtcdClient)

type ApplicationChangedNotification struct {
	Application string `json:"application"`
}

func notifyApp(appName string) {
	logrus.Infof("notifying registered clients for app %s", appName)
	conns := GetWatchAppConns(appName)
	if conns == nil {
		return
	}
	notification := ApplicationChangedNotification{Application:appName}
	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		logrus.Error("error during application changed notification", err)
	}
	for connId, conn := range conns {
		logrus.Debugf("notifying %s with %s", connId, jsonNotif)
		if err := conn.EmitMessage(jsonNotif); err != nil {
			logrus.Debug("error notifyng", err)
			RemoveWatchAppConn(appName, conn.ID())
		}
	}
}

func StartNotifier(appName string) {
	logrus.Infof("Start notifier for app %s", appName)
	if _, alreadyRun := notifiers[appName]; alreadyRun {
		logrus.Warnf("Notifier for app %s has been already started", appName)
	}

	etcdCl, err := etcd.New("root", cnf.GetString("etcd.password"))
	if err != nil {
		logrus.Errorf("Error starting notifier:%s", err)
		//wait and restart
		time.Sleep(2 * time.Second)
		go StartNotifier(appName)
	}
	defer etcdCl.Client.Close()
	watchCh := etcdCl.Client.Watch(context.Background(), appName + ".", clientv3.WithPrefix())
	notifiers[appName] = etcdCl
	for change := range watchCh {
		if change.Err() != nil {
			//restart notifier
			logrus.Error("Error during watching", err, ", restarting notifier")
			StopNotifier(appName)
			go StartNotifier(appName)
		} else {
			if !change.Canceled {
				logrus.Debugf("Change happened on etcd:%+v", change)
				notifyApp(appName)
			}
		}
	}
	logrus.Infof("Notifier for application %s stopped", appName)
}

func StopNotifier(appName string) {
	logrus.Infof("Stopping notifier for application %s", appName)
	watchEtcdCl, ok := notifiers[appName]
	if ok {
		watchEtcdCl.Client.Close()
	}
	delete(notifiers, appName)
}