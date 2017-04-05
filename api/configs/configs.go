package configs

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/api/utils"
	"github.com/foofilers/version"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /config api resources")
	configParty := router.Party("/configs", handlersFn...)
	configParty.Get("/:app/:version", getConfig)
	configParty.Put("/:app/:version", putConfig)
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
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()

	ver, err := version.Get(etcdCl, appName, appVersion);
	if utils.HandleError(ctx, err) {
		return
	}
	cnf, err := ver.GetConfig(etcdCl)
	if (utils.HandleError(ctx, err)) {
		return
	}
	ctx.JSON(iris.StatusOK, cnf)
}

func putConfig(ctx *iris.Context) {
	if (utils.MandatoryParams(ctx, "app", "version")) {
		return
	}
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	logrus.Infof("Put Configuration for app [%s] version:[%s]", appName, appVersion)
	appConfigs := make(map[string]string)
	if err := ctx.ReadJSON(&appConfigs); utils.HandleEtcdErrorMsg(ctx, err, " putConfig: Parsing JSON body") {
		return
	}
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	ver, err := version.Get(etcdCl, appName, appVersion);
	if utils.HandleError(ctx, err) {
		return
	}
	ver.OverwriteConfig(etcdCl, appConfigs)
	ctx.SetStatusCode(iris.StatusNoContent)
}