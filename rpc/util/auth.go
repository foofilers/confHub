package util

import (
	"google.golang.org/grpc/metadata"
	"golang.org/x/net/context"
	"errors"
	"github.com/foofilers/cfhd/auth"
	"google.golang.org/grpc/codes"
	"github.com/foofilers/cfhd/core/userManager"
	"github.com/sirupsen/logrus"
	"github.com/foofilers/cfhd/models"
	"google.golang.org/grpc/status"
	"strings"
)

func GetAuthUserId(ctx context.Context) (string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", errors.New("cannot retrieve metadata from context")
	}
	logrus.Debugf("metadata: %+v", md)
	authorizationHeaders := md["authorization"]
	logrus.Debugf("authorizationHeaders: %+v", authorizationHeaders)
	if len(authorizationHeaders) == 0 {

		return "", errors.New("authentication required")
	}
	bearer := strings.Split(authorizationHeaders[0], " ")
	if bearer[0] != "Bearer" {
		logrus.Errorf("authentication should be bearer")
		return "", errors.New("authentication should be bearer")
	}
	user, err := auth.ValidateJwt(bearer[1])
	if err != nil {
		logrus.Errorf("Error validating JWT:%s", err)
		return "", errors.New("authentication required")
	}
	return user, nil
}

func GetAuthUser(ctx context.Context) (*models.User, error) {
	var userId string
	var authError error
	if userId, authError = GetAuthUserId(ctx); authError != nil {
		logrus.Error(authError)
		return nil, status.Errorf(codes.Unauthenticated, authError.Error())
	}
	user, err := userManager.GetById(userId)
	if user == nil || err != nil {
		logrus.Errorf("Error retrieving user from db userId:%s err:%+v", userId, err)
		return nil, status.Errorf(codes.Unauthenticated, authError.Error())
	}
	return user, nil
}

