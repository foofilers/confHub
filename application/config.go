package application

import (
	"github.com/foofilers/confHub/etcd"
	"github.com/Sirupsen/logrus"
)

func (app *App) GetConfiguration(etcdCl *etcd.EtcdClient) (map[string]string, error) {
	values, err := etcdCl.GetWithPrefix(app.Name)
	if err != nil {
		return nil, err
	}
	if values == nil || len(values) == 0 {
		return nil, AppNotFoundError.Details(app.Name)
	}
	return values, nil
}

func (app *App) GetConfigurationVersion(etcdCl *etcd.EtcdClient, version string) (*Configuration, error) {
	var err error
	if len(version) == 0 {
		//use default version
		version, err = app.GetCurrentVersion(etcdCl)
		if err != nil {
			logrus.Errorf("Error getting current application version for app %s:%s", app.Name, err)
			return nil, err
		}
	}

	verExist, err := app.ExistVersion(version)
	if err != nil {
		return nil, err
	}
	if !verExist {
		return nil, VersionNotFound.Details(version, app.Name)
	}
	return &Configuration{app, version}, nil
}

