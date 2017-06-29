package users

import (
	"golang.org/x/net/context"
	"github.com/foofilers/cfhd/core/userManager"
	"github.com/foofilers/cfhd/models"
	"github.com/sirupsen/logrus"
	"github.com/golang/protobuf/ptypes/empty"
)

type UserService struct{}

func dbUser2User(dbUser *models.User) (*User) {
	if dbUser == nil {
		return nil
	}
	u := &User{}
	u.Id = dbUser.Id
	u.Username = dbUser.Username
	u.Email = dbUser.Email
	u.Permissions = make([]*Permission, len(dbUser.Permissions))
	for i, p := range (dbUser.Permissions) {
		u.Permissions[i] = &Permission{}
		u.Permissions[i].Application = p.Application
		u.Permissions[i].Perm = p.Perm
	}
	return u
}

func (usersRpc *UserService) List(request *UserListRequest, stream Users_ListServer) error {
	users, err := userManager.ListUsers(request.Query, request.Count, request.Page, request.Order)
	if err != nil {
		return err
	}
	for _, user := range (users) {
		if err := stream.Send(dbUser2User(&user)); err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func (usersRpc *UserService) Add(ctx context.Context, user *User) (*User, error) {
	dbUser := &models.User{}
	dbUser.Username = user.Username
	dbUser.Email = user.Email
	dbUser.Permissions = make([]models.Permission, len(user.Permissions))
	for i, perm := range user.Permissions {
		dbUser.Permissions[i].Application = perm.Application
		dbUser.Permissions[i].Perm = perm.Perm
	}
	err := userManager.AddUser(dbUser)
	if err != nil {
		return nil, err
	}
	user.Id = dbUser.Id;
	return user, nil;
}

func (usersRpc *UserService) Delete(ctx context.Context, req *DeleteRequest) (*empty.Empty, error) {
	return &empty.Empty{}, userManager.DeleteById(req.Id)
}

