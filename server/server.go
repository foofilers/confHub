package server

import (
	"github.com/foofilers/confHub/etcd"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	"github.com/foofilers/confHub/api"
	"golang.org/x/net/context"
	"github.com/foofilers/confHub/ws"
)

var app *iris.Framework

func Start(addr string) {
	StartAsync(addr, false)
}


func StartAsync(addr string, async bool) {
	etcd.InitEtcd()
	app = iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())
	app.Adapt(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:[]string{"GET","POST","PUT","DELETE"},
		AllowedHeaders:[]string{"*"},
		Debug:false,
	}))
	app.Adapt(view.HTML("./public", ".html"))
	ws.InitWs(app)



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