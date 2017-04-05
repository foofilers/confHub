package apps

import (
	"github.com/Sirupsen/logrus"
	iris "gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"github.com/foofilers/confHub/api/utils"
	"github.com/foofilers/confHub/application"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /apps api resources")
	appsParty := router.Party("/apps", handlersFn...)
	appsParty.Post("/", addApp)
	appsParty.Put("/:appName", updateApp)
	return appsParty
}

func addApp(ctx *iris.Context) {
	utils.MandatoryFormParams(ctx, "name")
	appName := ctx.FormValue("name")

	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()

	if _, err := application.Create(etcdCl, appName); utils.HandleError(ctx, err) {
		return
	}

	ctx.SetStatusCode(iris.StatusCreated)
}

func updateApp(ctx *iris.Context) {
	utils.MandatoryFormParams(ctx, "name")
	newName := ctx.FormValue("name")
	currentName := ctx.Param("appName")
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()
	app, err := application.Get(etcdCl, currentName)
	if utils.HandleError(ctx, err) {
		return
	}
	if app == nil {
		utils.HandleError(ctx,application.AppNotFoundError.Details(currentName))
		return
	}
	err = app.Rename(etcdCl, newName)
	if utils.HandleError(ctx, err) {
		return
	}
	ctx.SetStatusCode(iris.StatusNoContent)
}