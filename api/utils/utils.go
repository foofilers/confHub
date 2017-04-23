package utils

import (
	iris "gopkg.in/kataras/iris.v6"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/Sirupsen/logrus"
	"github.com/foofilers/confHub/application"
	"reflect"
)

func HandleError(ctx *iris.Context, err error) bool {
	return HandleErrorMsg(ctx, err, "%s")
}

func HandleErrorMsg(ctx *iris.Context, err error, format string, values ...interface{}) bool {
	if err == nil {
		return false
	}
	values = append(values, err)
	logrus.Debugf("Error class:%s", reflect.TypeOf(err))
	switch err {
	case rpctypes.ErrAuthFailed, rpctypes.ErrInvalidAuthToken, rpctypes.ErrPermissionDenied:
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusForbidden)
	case rpctypes.ErrRoleNotFound, application.ReferenceNotFoundError, application.TooManyReferenceLinksError:
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusPreconditionFailed)
	case rpctypes.ErrDuplicateKey:
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusConflict)
	case application.AppAlreadyExistError, application.VersionAlreadyExistError:
		logrus.Warnf(format, values...)
		ctx.EmitError(iris.StatusConflict)
	case application.AppNotFoundError, rpctypes.ErrUserNotFound,
		application.VersionNotFound, application.CurrentVersionNotSetted, application.ValueNotFoundError:
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
			logrus.Warnf("no parameter %s found", par)
			ctx.EmitError(iris.StatusPreconditionFailed)
			return true
		}
	}
	return false
}

func MandatoryFormParams(ctx *iris.Context, parameters ...string) bool {
	for _, par := range parameters {
		if len(ctx.FormValue(par)) == 0 {
			logrus.Warnf("no form value %s found", par)
			ctx.EmitError(iris.StatusPreconditionFailed)
			return true
		}
	}
	return false
}