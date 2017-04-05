package app

import (
	"github.com/Sirupsen/logrus"
	iris "gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/etcd"
	"github.com/foofilers/confHub/auth"
	"golang.org/x/net/context"
	"github.com/foofilers/confHub/api/utils"
	"time"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /apps api resources")
	appsParty := router.Party("/apps", handlersFn...)
	appsParty.Post("/", addApp)
	return appsParty
}

func addApp(ctx *iris.Context) {
	appName := ctx.FormValue("name")
	if len(appName) == 0 {
		ctx.EmitError(iris.StatusPreconditionFailed)
		ctx.Writef(" the 'name' parameter is mandatory")
		return
	}

	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if utils.HandleEtcdError(ctx, err) {
		return
	}
	defer etcdCl.Client.Close()

	key := appName + "._created"
	getResp, err := etcdCl.Client.Get(context.TODO(), key)
	if utils.HandleEtcdError(ctx, err) {
		return
	}
	if getResp.Count != 0 {
		logrus.Warn("Application already exist")
		ctx.EmitError(iris.StatusConflict)
		ctx.Writef(" the application %s already exists", appName)
		return
	}
	if _, err := etcdCl.Client.Put(context.TODO(), key, time.Now().Format(time.RFC3339)); utils.HandleEtcdError(ctx, err) {
		return
	}

	// creating roles
	for _, suffix := range []string{"RW", "R"} {
		if _, err := etcdCl.Client.RoleAdd(context.TODO(), appName + suffix); utils.HandleEtcdErrorMsg(ctx, err, "cannot create role %s", appName + suffix) {
			return
		}
	}
	ctx.SetStatusCode(iris.StatusCreated)
}