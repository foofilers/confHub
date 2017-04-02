package config

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/Sirupsen/logrus"
	"context"
	"github.com/foofilers/confHub/etcd"
	"github.com/coreos/etcd/client"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /config api resources")
	configParty := router.Party("/config", handlersFn...)
	configParty.Get("/:key", GetConfig)
	return configParty
}

func GetConfig(ctx *iris.Context) {
	confKey := ctx.Param("key")
	logrus.Infof("Get configuration: %s", confKey)
	etcdCl,_ := etcd.New("","")
	resp, err := etcdCl.Client.Get(context.Background(), confKey, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			ctx.NotFound()
		} else {
			ctx.Panic()
		}
	} else {
		ctx.Write(resp.Kvs[0].Value)
	}
}