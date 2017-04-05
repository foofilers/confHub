package application

import (
	"github.com/foofilers/confHub/etcd"
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

func (app *App) GetConfigurationVersion(etcdCl *etcd.EtcdClient, version string) (map[string]string, error) {
	values, err := etcdCl.GetWithPrefix(app.Name + "." + version)
	if err != nil {
		return nil, err
	}
	if values == nil || len(values) == 0 {
		return nil, AppNotFoundError.Details(app.Name)
	}
	return values, nil
}