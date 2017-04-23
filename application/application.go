package application

import (
	"github.com/foofilers/confHub/etcd"
	"context"
	"time"
	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
	"strings"
)

const CONFHUB_APPLICATION_NAMES_PREFIX = "confHub/applications/name/"
const CONFHUB_APPLICATION_VERSIONS_PREFIX = "confHub/applications/version/"

type App struct {
	Name           string `json:"name"`
	CurrentVersion string `json:"currentVersion"`
	CreatedAt      time.Time`json:"createdAt"`
}

func Exists(name string) (bool, error) {
	etcdCl, err := etcd.RootClient()
	if err != nil {
		return false, err
	}
	defer etcdCl.Client.Close()

	getResp, err := etcdCl.Client.Get(context.TODO(), CONFHUB_APPLICATION_NAMES_PREFIX + name, clientv3.WithCountOnly())
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	return getResp.Count != 0, nil
}

func Get(etcdCl*etcd.EtcdClient, name string) (*App, error) {
	exist, err := Exists(name)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, AppNotFoundError.Details(name)
	}

	key := name + "/_created"
	getResp, err := etcdCl.Client.Get(context.TODO(), key)
	if err != nil {
		logrus.Errorf("Error getting key:%s:%s", key, err)
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
	app.CurrentVersion, err = GetCurrentAppVersion(etcdCl, name)
	if err != nil && err != CurrentVersionNotSetted {
		logrus.Errorf("Error getting currentVersion for %s application:%s", name, err)
		return nil, err
	}
	return app, nil
}

func ListNames(etcdCl *etcd.EtcdClient) ([]string, error) {
	logrus.Info("Getting application list")
	getResp, err := etcdCl.Client.Get(context.TODO(), CONFHUB_APPLICATION_NAMES_PREFIX, clientv3.WithPrefix(), clientv3.WithKeysOnly(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	appNames := make([]string, 0)
	if err != nil {
		logrus.Error(err)
		return appNames, err
	}

	for _, k := range getResp.Kvs {
		appName := strings.Replace(string(k.Key), CONFHUB_APPLICATION_NAMES_PREFIX, "", 1)
		logrus.Debugf("Found app %s", appName)
		appNames = append(appNames, appName)
	}
	return appNames, nil
}

func List(etcdCl *etcd.EtcdClient) ([]*App, error) {
	logrus.Info("Getting application list")
	apps := make([]*App, 0)

	appNames, err := ListNames(etcdCl)
	if err != nil {
		logrus.Error(err)
		return apps, err
	}
	for _, appName := range appNames {
		app, err := Get(etcdCl, appName)
		if err != nil {
			return apps, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func Create(etcdCl *etcd.EtcdClient, name string) (*App, error) {
	exists, err := Exists(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, AppAlreadyExistError.Details(name)
	}

	ops := make([]clientv3.Op, 2, 2)

	app := &App{Name:name, CreatedAt:time.Now()}
	ops[0] = clientv3.OpPut(name + "/_created", app.CreatedAt.Format(time.RFC3339))
	ops[1] = clientv3.OpPut(CONFHUB_APPLICATION_NAMES_PREFIX + name, "true")

	if _, err := etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit(); err != nil {
		return nil, err
	}

	if err := createAppRoles(etcdCl, name); err != nil {
		return nil, err;
	}

	return app, nil
}

func (app *App) Rename(newName string) error {

	etcdCl, err := etcd.RootClient()
	if err != nil {
		return err
	}
	defer etcdCl.Client.Close()

	logrus.Infof("renaming application %s to %s", app.Name, newName)
	exists, err := Exists(newName)
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

	ops := make([]clientv3.Op, (len(appConf) * 2) + 2, (len(appConf) * 2) + 2)
	i := 0
	for k, v := range appConf {
		destKey := newName + "/" + k
		sourceKey := app.Name + "/" + k
		logrus.Debugf("move %s to %s", sourceKey, destKey)
		ops[i] = clientv3.OpPut(destKey, v)
		ops[i + 1] = clientv3.OpDelete(sourceKey)
		i += 2
	}
	//update applications list
	ops[i] = clientv3.OpPut(CONFHUB_APPLICATION_NAMES_PREFIX + newName, "true")
	i++
	ops[i] = clientv3.OpDelete(CONFHUB_APPLICATION_NAMES_PREFIX + app.Name)

	if _, err := etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit(); err != nil {
		return err
	}


	//regenerate grants
	if err := createAppRoles(etcdCl, newName); err != nil {
		return err;
	}

	//moving users to new role
	usersResp, err := etcdCl.Client.UserList(context.TODO())
	if err != nil {
		return err
	}
	for _, username := range usersResp.Users {
		userResp, err := etcdCl.Client.UserGet(context.TODO(), username)
		if err != nil {
			return err
		}
		for _, role := range userResp.Roles {
			if role == app.Name + "R" {
				logrus.Infof("add role %s to user %s", newName + "R", username)
				if _, err := etcdCl.Client.UserGrantRole(context.TODO(), username, newName + "R"); err != nil {
					return err
				}
			}
			if role == app.Name + "RW" {
				logrus.Infof("add role %s to user %s", newName + "RW", username)
				if _, err := etcdCl.Client.UserGrantRole(context.TODO(), username, newName + "RW"); err != nil {
					return err
				}
			}
		}
	}
	if err := app.removeAppRoles(etcdCl); err != nil {
		return err
	}

	app.Name = newName
	return nil
}

func createAppRoles(etcdCl *etcd.EtcdClient, appName string) error {
	// creating roles
	for _, suffix := range []string{"RW", "R"} {
		if _, err := etcdCl.Client.RoleAdd(context.TODO(), appName + suffix); err != nil {
			return err
		}
	}
	// associate roles to permission
	if _, err := etcdCl.Client.RoleGrantPermission(context.TODO(), appName + "RW", appName + "/", clientv3.GetPrefixRangeEnd(appName + "/"), clientv3.PermissionType(clientv3.PermReadWrite)); err != nil {
		return err
	}
	if _, err := etcdCl.Client.RoleGrantPermission(context.TODO(), appName + "R", appName + "/", clientv3.GetPrefixRangeEnd(appName + "/"), clientv3.PermissionType(clientv3.PermRead)); err != nil {
		return err
	}
	return nil
}

func (app *App) removeAppRoles(etcdCl *etcd.EtcdClient) error {
	//remove old permission
	logrus.Infof("delete role %s", app.Name + "R")
	if _, err := etcdCl.Client.RoleDelete(context.TODO(), app.Name + "R"); err != nil {
		return err
	}
	logrus.Infof("delete role %s", app.Name + "RW")
	if _, err := etcdCl.Client.RoleDelete(context.TODO(), app.Name + "RW"); err != nil {
		return err
	}
	return nil
}

func (app *App) Delete(etcdCl *etcd.EtcdClient) error {
	logrus.Infof("deleting application %s ", app.Name)
	appConf, err := app.GetConfiguration(etcdCl)
	if err != nil {
		return err
	}

	appVersions, err := app.GetVersions()
	if err != nil {
		return err
	}

	//delete keys
	ops := make([]clientv3.Op, 0)
	for k := range appConf {
		sourceKey := app.Name + "." + k
		logrus.Debugf("deleting %s", sourceKey)
		ops = append(ops, clientv3.OpDelete(sourceKey))
	}

	//delete app from app list
	ops = append(ops, clientv3.OpDelete(CONFHUB_APPLICATION_NAMES_PREFIX + app.Name))

	//delete versions
	for appVer := range appVersions {
		ops = append(ops, clientv3.OpDelete(CONFHUB_APPLICATION_VERSIONS_PREFIX + app.Name + appVer))
	}

	if _, err := etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit(); err != nil {
		logrus.Errorf("error deleting application %s:%s", app.Name, err)
		return err
	}

	if err := app.removeAppRoles(etcdCl); err != nil {
		return err
	}
	return nil
}