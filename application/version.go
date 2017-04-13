package application

import (
	"github.com/coreos/etcd/clientv3"
	"strings"
	"github.com/foofilers/confHub/etcd"
	"context"
	"github.com/Sirupsen/logrus"
)

func (app *App) GetVersions(etcdCl *etcd.EtcdClient) (map[string]bool, error) {
	versPrefix := CONFHUB_APPLICATIONS_PREFIX + app.Name + ".version."
	versResp, err := etcdCl.Client.Get(context.TODO(), versPrefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}
	res := make(map[string]bool)
	for _, kv := range versResp.Kvs {
		version := strings.Replace(string(kv.Key), versPrefix, "", 1)
		logrus.Debugf("found version %s in application %s", version, app.Name)
		res[version] = true
	}
	return res, nil
}

func (app *App) CreateVersion(etcdCl *etcd.EtcdClient, version string) error {
	currVers, err := app.GetVersions(etcdCl)
	if err != nil {
		return err
	}
	if _, found := currVers[version]; found {
		return VersionAlreadyExistError.Details(version, app.Name)
	}
	_, err = etcdCl.Client.Put(context.TODO(), CONFHUB_APPLICATIONS_PREFIX + app.Name + ".version." + version, "true")
	return err
}

func (app *App) ExistVersion(etcdCl *etcd.EtcdClient, version string) (bool, error) {
	currVers, err := app.GetVersions(etcdCl)
	if err != nil {
		return false, err
	}
	_, found := currVers[version];
	return found, nil
}

func (app *App) DeleteVersion(etcdCl *etcd.EtcdClient, version string) error {
	logrus.Debugf("Delete version:%+v", version)
	conf, err := app.GetConfigurationVersion(etcdCl, version)
	if err != nil {
		return err
	}
	currConfig, err := conf.GetConfig(etcdCl)
	if err != nil {
		return err
	}
	ops := make([]clientv3.Op, 0)
	for k := range currConfig {
		key := conf.GenKey(k)
		logrus.Debugf("deleting key:%s", key)
		ops = append(ops, clientv3.OpDelete(key))
	}
	ops = append(ops, clientv3.OpDelete(CONFHUB_APPLICATIONS_PREFIX + app.Name + ".version." + version))

	_, err = etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit()
	return err
}

func GetCurrentAppVersion(etcdCl *etcd.EtcdClient, appName string) (string, error) {
	currVerResp, err := etcdCl.Client.Get(context.TODO(),  appName + "._currentVersion")
	if err != nil {
		return "", err
	}
	if currVerResp.Count == 0 {
		return "", CurrentVersionNotSetted.Details(appName)
	}
	return string(currVerResp.Kvs[0].Value), nil
}

func (app *App) GetCurrentAppVersion(etcdCl *etcd.EtcdClient) (string, error) {
	return GetCurrentAppVersion(etcdCl, app.Name)
}

func (app *App) SetDefaultVersion(etcdCl *etcd.EtcdClient, version string) error {
	verExist, err := app.ExistVersion(etcdCl, version)
	if err != nil {
		return err
	}
	if !verExist {
		return VersionNotFound.Details(version, app.Name)
	}
	if _, err := etcdCl.Client.Put(context.TODO(),  app.Name + "._currentVersion", version); err != nil {
		return err
	}
	return nil
}