package versions

import (
	"github.com/Sirupsen/logrus"
	iris "gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/api/utils"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/application"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /config api resources")
	configParty := router.Party("/versions", handlersFn...)
	configParty.Post("/:app", addVersion)
	configParty.Get("/:app", getVersions)
	configParty.Put("/:app/:version", setDefaultVersion)
	configParty.Put("/:app/:version/copy", copyVersion)
	configParty.Delete("/:app/:version", deleteVersion)
	return configParty
}

func addVersion(ctx *iris.Context) {
	if (utils.MandatoryFormParams(ctx, "version")) {
		return
	}
	appName := ctx.Param("app")
	appVersion := ctx.FormValue("version")
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	app, err := application.Get(etcdCl, appName)
	if utils.HandleError(ctx, err) {
		return
	}
	if utils.HandleError(ctx, app.CreateVersion(etcdCl, appVersion)) {
		return
	}
	ctx.SetStatusCode(iris.StatusCreated)
}

func copyVersion(ctx *iris.Context) {
	if (utils.MandatoryFormParams(ctx, "version")) {
		return
	}
	appName := ctx.Param("app")
	srcVersion := ctx.Param("version")
	destVersion := ctx.FormValue("version")
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	app, err := application.Get(etcdCl, appName)
	if utils.HandleError(ctx, err) {
		return
	}

	cnf, err := app.GetConfigurationVersion(etcdCl, srcVersion)
	if utils.HandleError(ctx, err) {
		return
	}
	currCnf, err := cnf.GetConfig(etcdCl)
	if utils.HandleError(ctx, err) {
		return
	}

	if utils.HandleError(ctx, app.CreateVersion(etcdCl, destVersion)) {
		return
	}
	destCnf, err := app.GetConfigurationVersion(etcdCl, destVersion)
	if utils.HandleError(ctx, err) {
		return
	}

	if utils.HandleError(ctx, destCnf.SetConfig(etcdCl, currCnf)) {
		return
	}

	ctx.SetStatusCode(iris.StatusOK)
}

func getVersions(ctx *iris.Context) {
	appName := ctx.Param("app")
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	app, err := application.Get(etcdCl, appName)
	if utils.HandleError(ctx, err) {
		return
	}
	versions, err := app.GetVersions(etcdCl)
	if utils.HandleError(ctx, err) {
		return
	}
	res := make([]string, 0)
	for v := range versions {
		res = append(res, v)
	}
	ctx.JSON(iris.StatusOK, res)
}

func deleteVersion(ctx *iris.Context) {
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	app, err := application.Get(etcdCl, appName)
	if utils.HandleError(ctx, err) {
		return
	}
	if utils.HandleError(ctx, app.DeleteVersion(etcdCl, appVersion)) {
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent)
}

func setDefaultVersion(ctx *iris.Context) {
	appName := ctx.Param("app")
	appVersion := ctx.Param("version")
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	app, err := application.Get(etcdCl, appName)
	if utils.HandleError(ctx, err) {
		return
	}
	if utils.HandleError(ctx, app.SetDefaultVersion(etcdCl, appVersion)) {
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent)
}