package versions

import (
	"github.com/Sirupsen/logrus"
	iris "gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/api/utils"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/application"
	"github.com/foofilers/confHub/models"
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

	if utils.HandleError(ctx, app.CreateVersion(appVersion)) {
		return
	}

	// setting the default version if this version is the first one
	_, err = app.GetCurrentVersion(etcdCl)
	if err == application.CurrentVersionNotSetted {
		if utils.HandleError(ctx, app.SetDefaultVersion(etcdCl, appVersion)) {
			return
		}
	} else {
		if utils.HandleError(ctx, err) {
			return
		}
	}
	ctx.SetStatusCode(iris.StatusCreated)
}

func copyVersion(ctx *iris.Context) {
	if (utils.MandatoryFormParams(ctx, "version")) {
		return
	}
	appName := ctx.Param("app")
	srcVersion := ctx.Param("version")
	dstVersion := ctx.FormValue("version")
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

	if utils.HandleError(ctx, app.CreateVersion(dstVersion)) {
		return
	}
	destCnf, err := app.GetConfigurationVersion(etcdCl, dstVersion)
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
	versions, err := app.GetVersions()
	if utils.HandleError(ctx, err) {
		return
	}
	res := &models.ApplicationVersion{}
	res.Versions = make([]string, len(versions), len(versions))
	i := 0
	for v := range versions {
		res.Versions[i] = v
		i++
	}
	res.DefaultVersion, err = app.GetCurrentVersion(etcdCl)
	if err != application.CurrentVersionNotSetted && utils.HandleError(ctx, err) {
		return
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

	// check if it's the current version
	currVer, err := app.GetCurrentVersion(etcdCl)
	if err != application.CurrentVersionNotSetted && utils.HandleError(ctx, err) {
		return
	}
	if currVer == appVersion {
		logrus.Warn("cannot remove the current app version")
		ctx.SetStatusCode(iris.StatusPreconditionFailed)
		ctx.Writef("Cannot delete the current app version.")
		return
	}

	if utils.HandleError(ctx, app.DeleteVersion(appVersion)) {
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