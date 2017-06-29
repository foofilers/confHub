package applications

import (
	"golang.org/x/net/context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/foofilers/cfhd/models"
	"github.com/foofilers/cfhd/rpc/util"
	"github.com/foofilers/cfhd/core/applicationManager"
	"github.com/sirupsen/logrus"
)

type ApplicationService struct{}

func application2Db(app *Application) (*models.Application, error) {
	if app == nil {
		return nil, nil
	}
	configs, err := util.Struct2Map(app.Configuration)
	if err != nil {
		return nil, err
	}
	return &models.Application{
		Id:app.Id,
		Name:app.Name,
		Version:app.Version,
		Configuration:configs,
	}, nil
}

func db2Application(app *models.Application) (*Application, error) {
	if app == nil {
		return nil, nil
	}
	cfg, err := util.Map2Struct(app.Configuration)
	if err != nil {
		return nil, err
	}
	return &Application{
		Name:app.Name,
		Version:app.Version,
		Configuration:cfg,
	}, nil
}

func (service *ApplicationService) List(params *ApplicationListRequest, stream Applications_ListServer) error {
	logrus.Infof("ListApplication params:[%+v]", params)
	search, err := application2Db(params.Search)
	if err != nil {
		logrus.Error(err)
		return err
	}
	apps, err := applicationManager.ListApplication(search, params.Page, params.Count, params.Order)
	if err != nil {
		logrus.Error(err)
		return err
	}
	for _, app := range apps {
		grpcApp, err := db2Application(&app)
		if err != nil {
			logrus.Error(err)
			return err
		}
		stream.Send(grpcApp)
	}
	return nil
}

func (service *ApplicationService) Get(ctx context.Context, request *ApplicationGetRequest) (*Application, error) {
	app, err := applicationManager.GetApplication(request.Name, request.Version)
	if err != nil {
		return nil, err
	}
	return db2Application(app)
}

func (service *ApplicationService) Add(ctx context.Context, application *Application) (*Application, error) {
	app, err := application2Db(application)
	if err != nil {
		return nil, err
	}
	return application, applicationManager.AddApplication(app)
}

func (service *ApplicationService) Delete(ctx context.Context, request *DeleteRequest) (*empty.Empty, error) {

	return &empty.Empty{}, errors.New("Not implemented")
}

func (service *ApplicationService) Watch(params *ApplicationWatchRequest, stream Applications_WatchServer) error {
	/*var userId string
	var authError error
	if userId, authError = rpc.GetAuthUser(stream.Context()); authError != nil {
		return grpc.Errorf(code.Code_UNAUTHENTICATED, authError.Error())
	}*/

	appsCh := make(chan *models.Application)
	defer close(appsCh)
	stopCh := make(chan bool)
	defer close(stopCh)
	applicationManager.WatchApplication(params.Name, appsCh, stopCh)
	streamClosed := false
	for !streamClosed {
		select {
		case app := <-appsCh:
			rpcApp, err := db2Application(app)
			if err != nil {
				return err
			}
			err = stream.Send(&ApplicationWatch{Hearthbeat:rpcApp == nil, Application:rpcApp})
			if err != nil {
				streamClosed = true
				stopCh <- true
			}
		}
	}
	logrus.Debug("Watch terminated")
	return nil
}