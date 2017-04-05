package utils

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/Sirupsen/logrus"
)

func HandleEtcdError(ctx *iris.Context, err error) bool {
	return HandleEtcdErrorMsg(ctx, err, "%s")
}

func HandleEtcdErrorMsg(ctx *iris.Context, err error, format string, values ...interface{}) bool {
	if err == nil {
		return false
	}
	values = append(values, err)
	switch err {
	case rpctypes.ErrAuthFailed, rpctypes.ErrInvalidAuthToken :
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusForbidden)
	case rpctypes.ErrKeyNotFound, rpctypes.ErrUserNotFound :
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusForbidden)
	default:
		logrus.Errorf(format, values...)
		ctx.EmitError(iris.StatusInternalServerError)
	}
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