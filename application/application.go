package application

import (
	"github.com/foofilers/confHub/etcd"
	"context"
	"time"
	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
	"strings"
)

type App struct {
	Name      string
	CreatedAt time.Time
}

func Exists(etcdCl*etcd.EtcdClient, name string) (bool, error) {
	key := name + "._created"
	getResp, err := etcdCl.Client.Get(context.TODO(), key, clientv3.WithCountOnly())
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	return getResp.Count != 0, nil
}

func Get(etcdCl*etcd.EtcdClient, name string) (*App, error) {
	key := name + "._created"
	getResp, err := etcdCl.Client.Get(context.TODO(), key)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if getResp.Count == 0 {
		logrus.Warnf("No application found [key:%s]", key)
		return nil, nil
	}
	app := &App{Name:name}
	app.CreatedAt, err = time.Parse(time.RFC3339, string(getResp.Kvs[0].Value))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return app, nil
}

func List(etcdCl *etcd.EtcdClient) ([]*App, error) {
	logrus.Info("Getting application list")
	getResp, err := etcdCl.Client.Get(context.TODO(), "a", clientv3.WithFromKey())
	apps := make([]*App, 0)
	if err != nil {
		logrus.Error(err)
		return apps, err
	}
	appNames := make(map[string]bool)
	for _, k := range getResp.Kvs {
		appName := strings.Split(string(k.Key), ".")[0]
		appNames[appName] = true
	}

	for appName := range appNames {
		logrus.Debugf("Found app %s", appName)
		app, err := Get(etcdCl, appName)
		if err != nil {
			return apps, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func Create(etcdCl *etcd.EtcdClient, name string) (*App, error) {
	exists, err := Exists(etcdCl, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, AppAlreadyExistError.Details(name)
	}
	key := name + "._created"
	app := &App{Name:name, CreatedAt:time.Now()}
	if _, err := etcdCl.Client.Put(context.TODO(), key, app.CreatedAt.Format(time.RFC3339)); err != nil {
		return nil, err
	}

	// creating roles
	for _, suffix := range []string{"RW", "R"} {
		if _, err := etcdCl.Client.RoleAdd(context.TODO(), name + suffix); err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (app *App) Rename(etcdCl *etcd.EtcdClient, newName string) error {
	logrus.Infof("renaming application %s to %s", app.Name, newName)
	exists, err := Exists(etcdCl, newName)
	if err != nil {
		return err
	}
	if exists {
		return AppAlreadyExistError.Details(newName)
	}
	appConf, err := app.GetConfiguration(etcdCl)
	if err != nil {
		return err
	}

	ops := make([]clientv3.Op, len(appConf) * 2, len(appConf) * 2)
	i := 0
	for k, v := range appConf {
		destKey := newName + "." + k
		sourceKey := app.Name + "." + k
		logrus.Debugf("move %s to %s", sourceKey, destKey)
		ops[i] = clientv3.OpPut(destKey, v)
		ops[i + 1] = clientv3.OpDelete(sourceKey)
		i += 2
	}
	if _, err := etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit(); err != nil {
		return err
	}
	app.Name = newName
	return nil
}

func (app *App) Delete(etcdCl *etcd.EtcdClient) error {
	logrus.Infof("deleting application %s ", app.Name)
	appConf, err := app.GetConfiguration(etcdCl)
	if err != nil {
		return err
	}
	ops := make([]clientv3.Op, len(appConf), len(appConf))
	i := 0
	for k := range appConf {
		sourceKey := app.Name + "." + k
		logrus.Debugf("deleting %s", sourceKey)
		ops[i] = clientv3.OpDelete(sourceKey)
		i ++
	}
	if _, err := etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit(); err != nil {
		logrus.Errorf("error deleting application %s:%s", app.Name, err)
		return err
	}
	return nil
}