package values

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/api/utils"
	"github.com/foofilers/confHub/configuration"
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
	logrus.Infof("Get value appName:%s, version:%s, key:%s", appName, appVersion, confKey)
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()

	conf, err := configuration.Get(etcdCl, appName, appVersion)
	if utils.HandleError(ctx, err) {
		return
	}
	val, err := conf.GetValue(etcdCl, confKey)
	if utils.HandleError(ctx, err) {
		return
	}
	if val == nil {
		ctx.NotFound()
	}
	ctx.Write(val)
}

func putValue(ctx *iris.Context) {
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	confKey := ctx.Param("key")
	confValue := ctx.FormValue("value")
	logrus.Infof("Put value appName:%s, version:%s, key:%s, value:%s", appName, appVersion, confKey, confValue)

	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()

	conf, err := configuration.Get(etcdCl, appName, appVersion)
	if utils.HandleError(ctx, err) {
		return
	}
	err = conf.PutValue(etcdCl, confKey, confValue)
	if utils.HandleError(ctx, err) {
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent);
}

func deleteValue(ctx *iris.Context) {
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	confKey := ctx.Param("key")
	logrus.Infof("Delete value appName:%s, version:%s, key:%s", appName, appVersion, confKey)
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()

	conf, err := configuration.Get(etcdCl, appName, appVersion)
	if utils.HandleError(ctx, err) {
		return
	}
	err = conf.DeleteValue(etcdCl, confKey)
	if utils.HandleError(ctx, err) {
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent);
}