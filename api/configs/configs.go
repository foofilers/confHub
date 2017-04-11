package configs

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/api/utils"
	"github.com/foofilers/confHub/application"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /config api resources")
	configParty := router.Party("/configs", handlersFn...)
	configParty.Get("/:app/:version", getConfig)
	configParty.Get("/:app", getConfig)
	configParty.Put("/:app/:version", putConfig)
	return configParty
}

func getConfig(ctx *iris.Context) {
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")

	logrus.Infof("Get configuration for app: %s version:%s", appName, appVersion)
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()

	app,err:=application.Get(etcdCl,appName)
	if utils.HandleError(ctx, err) {
		return
	}

	appConf, err := app.GetConfigurationVersion(etcdCl, appVersion);
	if utils.HandleError(ctx, err) {
		return
	}
	cnf, err := appConf.GetConfig(etcdCl)
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
	defer etcdCl.Client.Close()

	app,err:=application.Get(etcdCl,appName)
	if utils.HandleError(ctx, err) {
		return
	}

	appConf, err := app.GetConfigurationVersion(etcdCl, appVersion);
	if utils.HandleError(ctx, err) {
		return
	}
	if utils.HandleError(ctx, appConf.SetConfig(etcdCl, appConfigs)) {
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent)
}
