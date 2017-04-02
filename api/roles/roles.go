package roles

import (
	"github.com/Sirupsen/logrus"
	iris "gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/etcd"
	"golang.org/x/net/context"
	"github.com/foofilers/confHub/auth"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /roles api resources")
	usersParty := router.Party("/roles", handlersFn...)
	usersParty.Get("/", getRoles)
	return usersParty
}

func getRoles(ctx *iris.Context) {
	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	etcdCl.Client.RoleAdd(context.TODO(),"test2")
	roles, err := etcdCl.Client.RoleList(context.TODO())
	if err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	logrus.Debug("roles resp %+v",roles)
	for _, role := range roles.Roles {
		logrus.Debugf("getting role %s", role)
		roleInfo, err := etcdCl.Client.RoleGet(context.TODO(), role)
		if err != nil {
			logrus.Error(err)
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}
		logrus.Debugf("roleInfo:%+v", roleInfo)
	}
	ctx.WriteString("ok")
}
