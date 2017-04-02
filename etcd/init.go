package etcd

import (
	log "github.com/Sirupsen/logrus"
	cnf "github.com/spf13/viper"
	"context"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

func InitEtcd() {
	rootCl, err := New(cnf.GetString("etcd.username"), cnf.GetString("etcd.password"))
	if err != nil {
		if err == rpctypes.ErrAuthFailed {
			//try to enable auth
			initEtcdAuthentication()
			if rootCl, err = New(cnf.GetString("etcd.username"), cnf.GetString("etcd.password")); err != nil {
				log.Fatal(err)
			}
		}
		log.Fatal(err)
	}
	defer rootCl.Client.Close()
	log.Debug("Checking confHub settings in etcd system")
	if resp, err := rootCl.Client.Get(context.TODO(), "confHub.version"); err != nil {
		log.Fatal(err)
	} else {
		if len(resp.Kvs) == 0 {
			createBasicConfiguration(rootCl)
		} else {
			log.Debugf("check schema version etcd Version:%s current:%s", resp.Kvs[0].Value, cnf.GetString("version"))
			if string(resp.Kvs[0].Value) != cnf.GetString("version") {
				//different confHub version
				updateSchema(string(resp.Kvs[0].Value))
			}
		}
	}
}

func initEtcdAuthentication() {
	log.Info("Enabling web authentication")
	rootCl, err := New("", "");
	if err != nil {
		log.Fatal(err)
	}
	if _, err = rootCl.Client.UserAdd(context.TODO(), "root", "confHub"); err != nil {
		log.Fatal(err)
	}

	log.Infof("root user created with default password 'confHub', PLEASE CHANGE IT")

	if _, err = rootCl.Client.UserGrantRole(context.TODO(), "root", "root"); err != nil {
		log.Fatal(err)
	}

	_, err = rootCl.Client.Auth.AuthEnable(context.TODO())
	if (err != nil) {
		log.Fatal(err)
	}

	rootCl, err = New(cnf.GetString("etcd.username"), cnf.GetString("etcd.password"))
	if err != nil {
		log.Fatal(err)
	}
}

/**
	Updating etcd schemas
 */
func updateSchema(fromVersion string) {
	log.Infof("Updating confHub Etcd schema from version %s to %s", fromVersion, cnf.GetString("version"))
}

func createBasicConfiguration(rootCl *EtcdClient) {
	log.Info("Creating confhub configuration entries")
	resp, err := rootCl.Client.Put(context.Background(), "confHub.version", "0.0.1b")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Debug(resp)
	}
}
