package etcd

import (
	"github.com/coreos/etcd/clientv3"
	log "github.com/Sirupsen/logrus"
	cnf "github.com/spf13/viper"
	"time"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/foofilers/confHub/auth"
	"github.com/coreos/pkg/cryptoutil"
	"encoding/base64"
)

type EtcdClient struct {
	Client *clientv3.Client
}

func New(username, password string) (*EtcdClient, error) {
	servers := cnf.GetStringSlice("etcd.servers")
	log.Infof("Initializing Etcd Client userId:%v, password:%v servers:%+v", username, password, servers)
	tlsCfg := transport.TLSInfo{
		ClientCertAuth:true,
	}
	tlsConfig, err := tlsCfg.ClientConfig()
	tlsConfig.InsecureSkipVerify = true
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cfg := clientv3.Config{
		Endpoints:              servers,
		DialTimeout: 5 * time.Second,
		Username: username,
		Password: password,
		TLS:tlsConfig,
	}
	var clErr error
	cl := &EtcdClient{}
	cl.Client, clErr = clientv3.New(cfg)
	if clErr != nil {
		return nil, clErr
	}

	return cl, nil
}

func LoggedClient(user auth.LoggedUser) (*EtcdClient, error) {
	unbase64Pwd, _ := base64.StdEncoding.DecodeString(user.CryptedPassword)
	clearedPassword, err := cryptoutil.AESDecrypt(unbase64Pwd, []byte(cnf.GetString("pwdSecretKey")))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return New(user.Username, string(clearedPassword))
}
