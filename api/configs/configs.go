package configs

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/Sirupsen/logrus"
	"context"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/coreos/etcd/clientv3"
	"strings"
	"github.com/foofilers/confHub/api/utils"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /config api resources")
	configParty := router.Party("/configs", handlersFn...)
	configParty.Get("/:app/:version", GetConfig)
	configParty.Put("/:app/:version", PutConfig)
	return configParty
}

func getConfig(ctx *iris.Context) {
	if (utils.MandatoryParams(ctx, "app", "version")) {
		return
	}
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	logrus.Infof("Get configuration for app: %s version:%s", appName, appVersion)
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleEtcdError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()

	keyPrefix := appName + "." + appVersion
	resp, err := etcdCl.Client.Get(context.TODO(), keyPrefix, clientv3.WithPrefix())
	if utils.HandleEtcdError(ctx, err) {
		return
	}
	if resp.Count == 0 {
		ctx.NotFound()
		return
	}
	confMap := make(map[string]string)
	for _, kv := range resp.Kvs {
		logrus.Debugf("%+v", kv)
		fullKey := string(kv.Key)
		confMap[strings.Replace(fullKey, keyPrefix + ".", "", 1)] = string(kv.Value)
	}
	ctx.JSON(iris.StatusOK, confMap)
}

func putConfig(ctx *iris.Context) {
	if (utils.MandatoryParams(ctx, "app", "version")) {
		return
	}
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	appConfigs := make(map[string]string)
	if err := ctx.ReadJSON(appConfigs); utils.HandleEtcdErrorMsg(ctx, err, " putConfig: Parsing JSON body") {
		return
	}
	logrus.Infof("Get configuration for app: %s version:%s", appName, appVersion)
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleEtcdError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()
	txn := etcdCl.Client.Txn(context.TODO())
	ops := make([]clientv3.Op, len(appConfigs), len(appConfigs))
	for k, v := range appConfigs {
		ops = append(ops, clientv3.OpPut(k, v))
	}
	_, err = txn.Then(ops).Commit()
	if utils.HandleEtcdError(ctx, err) {
		return
	}
}