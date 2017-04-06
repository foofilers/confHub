package configuration

import (
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/application"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"github.com/Sirupsen/logrus"
)

type Configuration struct {
	Application *application.App
	Version     string
}

func New(app *application.App, version string) *Configuration {
	return &Configuration{app, version}
}

func Get(etcdCl *etcd.EtcdClient, appName, version string) (*Configuration, error) {
	app, err := application.Get(etcdCl, appName)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, application.AppNotFoundError.Details(appName)
	}
	return New(app, version), nil
}

func (version *Configuration) GetConfig(etcdCl *etcd.EtcdClient) (map[string]string, error) {
	return etcdCl.GetWithPrefix(version.Application.Name + "." + version.Version)
}

func (version *Configuration) genKey(key string) string {
	return version.Application.Name + "." + version.Version + "." + key
}

func (version *Configuration) SetConfig(etcdCl *etcd.EtcdClient, newConf map[string]string) error {
	logrus.Debugf("SetConfig for version:%+v with values:%+v", version, newConf)
	currConfig, err := version.GetConfig(etcdCl)
	if err != nil {
		return err
	}
	ops := make([]clientv3.Op, 0)
	for k := range currConfig {
		key := version.genKey(k)
		logrus.Debugf("deleting key:%s", key)
		ops = append(ops, clientv3.OpDelete(key))
	}
	for k, v := range newConf {
		key := version.genKey(k)
		logrus.Debugf("insert key:%s=%s", key, v)
		ops = append(ops, clientv3.OpPut(version.genKey(k), v))
	}
	_, err = etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit()
	return err
}

func (version *Configuration) Delete(etcdCl *etcd.EtcdClient) error {
	logrus.Debugf("Delete version:%+v", version)
	currConfig, err := version.GetConfig(etcdCl)
	if err != nil {
		return err
	}
	ops := make([]clientv3.Op, 0)
	for k := range currConfig {
		key := version.genKey(k)
		logrus.Debugf("deleting key:%s", key)
		ops = append(ops, clientv3.OpDelete(key))
	}
	_, err = etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit()
	return err
}