package rpc

import (
	"net"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"github.com/foofilers/cfhd/rpc/users"
	"github.com/foofilers/cfhd/rpc/applications"
	"github.com/sirupsen/logrus"
	"sync"
	auth_rpc "github.com/foofilers/cfhd/rpc/auth"
)

var grpcServer *grpc.Server



func Start(port string, wg *sync.WaitGroup, quitCh chan bool) {
	logrus.Info("starting GRPC")
	wg.Add(1)
	go func() {
		select {
		case <-quitCh:
			Stop()
			wg.Done()
		}
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	grpcServer = grpc.NewServer()

	auth_rpc.RegisterAuthServer(grpcServer, &auth_rpc.AuthService{})
	users.RegisterUsersServer(grpcServer, &users.UserService{})
	applications.RegisterApplicationsServer(grpcServer, &applications.ApplicationService{})

	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Errorf("failed to serve: %v", err)
		}
	}()
}

func Stop() {
	grpc.NewServer()
	logrus.Info("Stopping grpc server")
	grpcServer.GracefulStop()
	logrus.Info("grpc server stopped")
}
