package main

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"github.com/foofilers/confHub/api"
	"github.com/foofilers/confHub/conf"
	 "github.com/foofilers/confHub/etcd"
	 "github.com/foofilers/confHub/log"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	"os"
)

func main() {
	os.Setenv("ETCDCTL_API","3")
	conf.InitConf()
	log.InitLog()
	etcd.InitEtcd()
	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())
	app.Adapt(cors.Default())
	app.Adapt(view.HTML("./public",".html"))

	app.Get("/", func(ctx *iris.Context) {
		ctx.MustRender("index.html",nil)
	})

	initRoutes(app)

	app.Listen("0.0.0.0:8080")

}

func initRoutes(app *iris.Framework) {
	api.InitApi(app.Party("/api"))
}