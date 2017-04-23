package application

import (
	"github.com/foofilers/confHub/etcd"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"github.com/Sirupsen/logrus"
	"strings"
	"github.com/spf13/viper"
)

type Configuration struct {
	Application *App
	Version     string
}

func getMaxReferences() int {
	if viper.IsSet("maxReference") {
		return viper.GetInt("maxReference")
	} else {
		return 10
	}
}

func (version *Configuration) GetConfig(etcdCl *etcd.EtcdClient) (map[string]string, error) {
	return etcdCl.GetWithPrefix(version.Application.Name + "/" + version.Version)
}

func (version *Configuration) GetConfigFollowingReferences(etcdCl *etcd.EtcdClient) (map[string]string, error) {
	rawConfig, err := version.GetConfig(etcdCl)
	if err != nil {
		return rawConfig, err
	}
	result := make(map[string]string)
	for k, v := range rawConfig {
		result[k], err = version.getRefValue(etcdCl, v, getMaxReferences())
		if err != nil {
			return result, err
		}
	}
	return result, nil;
}

func (version *Configuration) getRefValue(etcdCl *etcd.EtcdClient, value string, maxRefs int) (string, error) {
	if maxRefs <= 0 {
		return "", TooManyReferenceLinksError.Details(version.Application.Name, version.Version)
	}
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		//link reference
		ref := value[2:len(value) - 1]
		refValueResp, err := etcdCl.Client.Get(context.TODO(), ref)
		if err != nil {
			logrus.Debugf("Error getting ref [%s]:%s", ref, err)
			return "", err
		}
		if refValueResp.Count == 0 {
			return "", ReferenceNotFoundError.Details(ref, version.Application.Name)
		}
		refValue := string(refValueResp.Kvs[0].Value)
		return version.getRefValue(etcdCl, refValue, maxRefs - 1)
	} else {
		return value, nil
	}
}

func (version *Configuration) GenKey(key string) string {
	return version.Application.Name + "/" + version.Version + "/" + key
}

func (version *Configuration) SetConfig(etcdCl *etcd.EtcdClient, newConf map[string]string) error {
	logrus.Debugf("SetConfig for version:%+v with values:%+v", version, newConf)
	currConfig, err := version.GetConfig(etcdCl)
	if err != nil {
		return err
	}
	ops := make([]clientv3.Op, 0)
	for k := range currConfig {
		key := version.GenKey(k)
		logrus.Debugf("deleting key:%s", key)
		ops = append(ops, clientv3.OpDelete(key))
	}
	for k, v := range newConf {
		key := version.GenKey(k)
		logrus.Debugf("insert key:%s=%s", key, v)
		ops = append(ops, clientv3.OpPut(version.GenKey(k), v))
	}
	_, err = etcdCl.Client.Txn(context.TODO()).Then(ops...).Commit()
	return err
}

func (version *Configuration) GetValue(etcdCl *etcd.EtcdClient, key string) (string, error) {
	fullKey := version.Application.Name + "/" + version.Version + "/" + key
	resp, err := etcdCl.Client.Get(context.TODO(), fullKey)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", ValueNotFoundError.Details(key, version.Application.Name, version.Version)
	}
	return string(resp.Kvs[0].Value), nil
}

func (version *Configuration) GetValueFollowingReference(etcdCl *etcd.EtcdClient, key string) (string, error) {
	value, err := version.GetValue(etcdCl, key)
	if err != nil {
		logrus.Errorf("Error getting value %s/%s/%s:%s", version.Application.Name, version.Version, key, err)
		return "", err
	}
	return version.getRefValue(etcdCl, value, getMaxReferences())
}

func (version *Configuration) DeleteValue(etcdCl *etcd.EtcdClient, key string) error {
	fullKey := version.Application.Name + "/" + version.Version + "/" + key
	if _, err := etcdCl.Client.Delete(context.TODO(), fullKey); err != nil {
		return err
	}
	return nil
}

func (version *Configuration) RenameAndSetValue(etcdCl *etcd.EtcdClient, key, newKey, value string) error {
	delKey := version.Application.Name + "/" + version.Version + "/" + key
	fullKey := version.Application.Name + "/" + version.Version + "/" + newKey
	logrus.Debugf("RenameAndSetConfig from %s -> %s  value:%s", delKey, fullKey, value)
	_, err := etcdCl.Client.Txn(context.TODO()).Then(clientv3.OpDelete(delKey), clientv3.OpPut(fullKey, value)).Commit()
	return err;
}

func (version *Configuration) PutValue(etcdCl *etcd.EtcdClient, key, value string) error {
	fullKey := version.Application.Name + "/" + version.Version + "/" + key
	if _, err := etcdCl.Client.Put(context.TODO(), fullKey, value); err != nil {
		return err
	}
	return nil
}

