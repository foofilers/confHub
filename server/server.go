package server

import (
	"github.com/foofilers/confHub/conf"
	"github.com/foofilers/confHub/log"
	"github.com/foofilers/confHub/etcd"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	"github.com/foofilers/confHub/api"
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

var app *iris.Framework

func Start(addr string) {
	StartAsync(addr, false)
}

func StartAsync(addr string, async bool) {
	if !conf.IsInitialized() {
		logrus.Fatal("Configuration is not initialized, pleas run conf.InitConfFromFile or conf.conf.InitConf")
	}
	log.InitLog()
	etcd.InitEtcd()
	app = iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())
	app.Adapt(cors.Default())
	app.Adapt(view.HTML("./public", ".html"))

	app.Get("/", func(ctx *iris.Context) {
		ctx.MustRender("index.html", nil)
	})

	api.InitApi(app.Party("/api"))
	if (async) {
		go app.Listen(addr)
	} else {
		app.Listen(addr)
	}
}

func Stop() {
	if app != nil {
		app.Shutdown(context.TODO());
	}
}