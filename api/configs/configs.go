package configs

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/api/utils"
	"github.com/foofilers/configuration"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /config api resources")
	configParty := router.Party("/configs", handlersFn...)
	configParty.Get("/:app/:version", getConfig)
	configParty.Put("/:app/:version", putConfig)
	configParty.Delete("/:app/:version", deleteConfig)
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

	appConf, err := configuration.Get(etcdCl, appName, appVersion);
	if utils.HandleError(ctx, err) {
		return
	}
	cnf, err := appConf.GetConfig(etcdCl)
	if (utils.HandleError(ctx, err)) {
		return
	}
	if len(cnf) == 0 {
		ctx.NotFound()
	} else {
		ctx.JSON(iris.StatusOK, cnf)
	}
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
	appConf, err := configuration.Get(etcdCl, appName, appVersion);
	if utils.HandleError(ctx, err) {
		return
	}
	if utils.HandleError(ctx, appConf.SetConfig(etcdCl, appConfigs)) {
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent)
}

func deleteConfig(ctx *iris.Context) {
	if (utils.MandatoryParams(ctx, "app", "version")) {
		return
	}
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	appConf, err := configuration.Get(etcdCl, appName, appVersion);
	if utils.HandleError(ctx, err) {
		return
	}
	if utils.HandleError(ctx, appConf.Delete(etcdCl)) {
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent)
}