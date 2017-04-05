package users

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/foofilers/confHub/models"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/etcd"
	"golang.org/x/net/context"
	"github.com/foofilers/confHub/auth"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

func InitAPI(router *iris.Router, handlersFn ...iris.HandlerFunc) *iris.Router {
	logrus.Info("initializing /users api resources")
	usersParty := router.Party("/users", handlersFn...)
	usersParty.Post("/", addUser)
	return usersParty
}

func addUser(ctx *iris.Context) {
	user := &models.User{}
	if err := ctx.ReadJSON(user); err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	logrus.Debugf("adding user: %+v", user)

	etcdCl, err := etcd.LoggedClient(ctx.Get("LoggedUser").(auth.LoggedUser))
	if err != nil {
		logrus.Error(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	//check role presence
	for _, role := range user.Roles {
		if _, err := etcdCl.Client.RoleGet(context.TODO(), role); err != nil {
			if err == rpctypes.ErrRoleNotFound {
				ctx.EmitError(iris.StatusPreconditionFailed)
				ctx.Writef(" No role found with name:%s", role)
				return
			} else {
				logrus.Error(err)
				ctx.EmitError(iris.StatusInternalServerError)
				return
			}
		}
	}

	// adding user
	if _, err = etcdCl.Client.UserAdd(context.TODO(), user.Username, user.Password); err != nil {
		if err == rpctypes.ErrUserAlreadyExist {
			ctx.EmitError(iris.StatusPreconditionFailed)
			ctx.Writef(" User %s already exists", user.Username)
			return
		} else {
			logrus.Error(err)
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}
	}
	// setting user roles
	for _, role := range user.Roles {
		if _, err = etcdCl.Client.UserGrantRole(context.TODO(), user.Username, role); err != nil {
			logrus.Error(err)
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}
	}
	ctx.SetStatusCode(iris.StatusCreated)

}
