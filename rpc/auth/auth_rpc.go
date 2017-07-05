package auth

import (
	"golang.org/x/net/context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/foofilers/cfhd/core/userManager"
	auth_core"github.com/foofilers/cfhd/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type AuthService struct{}

func (service *AuthService) Login(ctx context.Context, req *LoginRequest) (*Jwt, error) {
	user, err := userManager.Login(req.Username, req.Password)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Invalid credentials")
	}

	//todo manage expirationTime
	jwt, err := auth_core.GenerateJwt(user.Id.Hex(), -1)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Invalid credentials")
	}
	return &Jwt{Jwt:jwt}, nil
}

func (service *AuthService) Logout(ctx context.Context, in *LogoutRequest) (*empty.Empty, error) {
	//todo manage logout
	return &empty.Empty{}, nil
}
