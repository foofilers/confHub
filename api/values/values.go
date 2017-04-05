package values

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/Sirupsen/logrus"
	"context"
	"github.com/foofilers/confHub/etcd"
	"github.com/coreos/etcd/client"
	"github.com/foofilers/confHub/auth"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /config api resources")
	configParty := router.Party("/values", handlersFn...)
	configParty.Get("/:app/:version/:key", getValue)
	configParty.Put("/:app/:version/:key", putValue)
	configParty.Delete("/:app/:version/:key", deleteValue)
	return configParty
}

func getValue(ctx *iris.Context) {
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	confKey := ctx.Param("key")
	logrus.Infof("Get configuration: %s", confKey)
	etcdCl, _ := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	defer etcdCl.Client.Close()

	fullKey := appName + "." + appVersion + "." + confKey
	resp, err := etcdCl.Client.Get(context.TODO(), fullKey, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			ctx.NotFound()
		} else {
			ctx.Panic()
		}
	} else {
		ctx.Write(resp.Kvs[0].Value)
	}
}

func putValue(ctx *iris.Context) {
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	confKey := ctx.Param("key")
	confValue := ctx.FormValue("value")

	// TODO check if the app and version exists

	logrus.Infof("Get configuration: %s", confKey)
	etcdCl, _ := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	defer etcdCl.Client.Close()

	fullKey := appName + "." + appVersion + "." + confKey
	if _, err := etcdCl.Client.Put(context.TODO(), fullKey, confValue); err != nil {
		logrus.Errorf("error adding/updating key:%s error:%s", fullKey, err.Error())
		ctx.EmitError(iris.StatusInternalServerError)
	}
}

func deleteValue(ctx *iris.Context) {
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	confKey := ctx.Param("key")
	logrus.Infof("Get configuration: %s", confKey)
	etcdCl, _ := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	defer etcdCl.Client.Close()

	fullKey := appName + "." + appVersion + "." + confKey
	if _, err := etcdCl.Client.Delete(context.TODO(), fullKey); err != nil {
		logrus.Errorf("error deleting key:%s error:%s", fullKey, err.Error())
		ctx.EmitError(iris.StatusInternalServerError)
	}
}