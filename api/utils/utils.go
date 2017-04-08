package utils

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/application"
	"reflect"
)

func HandleError(ctx *iris.Context, err error) bool {
	return HandleEtcdErrorMsg(ctx, err, "%s")
}

func HandleEtcdErrorMsg(ctx *iris.Context, err error, format string, values ...interface{}) bool {
	if err == nil {
		return false
	}
	values = append(values, err)
	logrus.Debugf("Error class:%s",reflect.TypeOf(err))
	switch err {
	case rpctypes.ErrAuthFailed, rpctypes.ErrInvalidAuthToken :
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusForbidden)
	case rpctypes.ErrRoleNotFound:
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusPreconditionFailed)
	case rpctypes.ErrDuplicateKey:
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusConflict)
	case application.AppAlreadyExistError:
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusConflict)
	case application.AppNotFoundError, rpctypes.ErrUserNotFound:
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusNotFound)
	default:
		logrus.Errorf(format, values...)
		ctx.EmitError(iris.StatusInternalServerError)
	}
	ctx.Writef(format, values...)
	return true
}

func MandatoryParams(ctx *iris.Context, parameters ...string) bool {
	for _, par := range parameters {
		if len(ctx.Param(par)) == 0 && len(ctx.FormValue(par)) == 0 {
			ctx.EmitError(iris.StatusPreconditionFailed)
			return true
		}
	}
	return false
}

func MandatoryFormParams(ctx *iris.Context, parameters ...string) bool {
	for _, par := range parameters {
		if len(ctx.FormValue(par)) == 0 {
			ctx.EmitError(iris.StatusPreconditionFailed)
			return true
		}
	}
	return false
}