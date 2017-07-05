package applications

import (
	"testing"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/foofilers/cfhd/util"
	rpcUtil "github.com/foofilers/cfhd/rpc/util"
	"reflect"
	"io"
	"time"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/foofilers/cfhd/rpc/auth"
	"google.golang.org/grpc/metadata"
	"log"
	"github.com/sirupsen/logrus"
)

var Conn *grpc.ClientConn
var createdApplication *Application

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	var err error
	Conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer Conn.Close()
	if err != nil {
		panic(err)
	}

	m.Run()
}

func getAuthContext(username, password string) context.Context {
	authClient := auth.NewAuthClient(Conn)
	jwt, err := authClient.Login(context.TODO(), &auth.LoginRequest{
		Username:username,
		Password:password,
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return metadata.NewContext(context.TODO(), metadata.Pairs("authorization", "Bearer " + jwt.Jwt))
}

func TestAddApplication(t *testing.T) {

	applClient := NewApplicationsClient(Conn)
	app := &Application{
		Name:util.RandStringRunes(8),
		Version:"1.0.0",
	}
	var err error
	app.Configuration, err = rpcUtil.Map2Struct(map[string]interface{}{
		"strProp1":"val1",
		"intProp1":1,
		"floatProp1":1.5,
		"subStruct":struct {
			StructPropA string
			StructPropB int
		}{
			StructPropA:"ciao",
			StructPropB:1,
		},
		"map":map[string]interface{}{
			"subMap1":"submap1Val",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", app.Configuration)
	_, err = applClient.Add(getAuthContext("admin","admin"), app)
	if err != nil {
		t.Fatal(err)
	}
	createdApplication = app
}

func TestGetApplication(t *testing.T) {
	if createdApplication == nil {
		t.Skip("previous test [TestAddApplication] in error ")
	}
	applClient := NewApplicationsClient(Conn)
	getApp, err := applClient.Get(context.TODO(), &ApplicationGetRequest{Name:createdApplication.Name, Version:createdApplication.Version})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(getApp, createdApplication) {
		t.Fatal("not the same application")
	}
	t.Logf("application:%+v", getApp)
}

func TestListAllApplication(t *testing.T) {
	if createdApplication == nil {
		t.Skip("previous test [TestAddApplication] in error ")
	}
	applClient := NewApplicationsClient(Conn)

	apps, err := applClient.List(getAuthContext("admin","admin"), &ApplicationListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for {
		app, err := apps.Recv()
		if err == io.EOF {
			break;
		}
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(app, createdApplication) {
			return
		}
	}
	t.Fatal("cannot find created application on list")
}

func TestListSpecificApplication(t *testing.T) {
	var err error
	if createdApplication == nil {
		t.Skip("previous test [TestAddApplication] in error ")
	}
	applClient := NewApplicationsClient(Conn)
	apps, err := applClient.List(getAuthContext("admin","admin"), &ApplicationListRequest{
		Search:&Application{
			Name:createdApplication.Name,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	i := 0
	var app *Application
	for {
		currApp, err := apps.Recv()
		if err == io.EOF {
			break;
		}
		if err != nil {
			t.Fatal(err)
		}
		app = currApp
		logrus.Debugf("%+v", currApp)
		i++
	}
	if i != 1 {
		t.Fatalf("Should be found only one application but I've found %d", i)
	}
	if !reflect.DeepEqual(app, createdApplication) {
		t.Errorf("%+v\nshould be\n%+v", app, createdApplication)
		t.Fatal("not the same application")
	}
}

func TestWatch(t *testing.T) {
	var err error
	if createdApplication == nil {
		t.Skip("previous test [TestAddApplication] in error ")
	}
	applClient := NewApplicationsClient(Conn)

	apps, err := applClient.Watch(getAuthContext("admin","admin"), &ApplicationWatchRequest{Name:createdApplication.Name})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for i := 0; i < 5; i++ {
			var newApp Application
			copier.Copy(&newApp, createdApplication)
			newApp.Version = fmt.Sprintf("2.0.%d", i)
			applClient.Add(getAuthContext("admin","admin"), &newApp)
			time.Sleep(2 * time.Second)
		}
	}()
	i := 0;
	for {
		currApp, err := apps.Recv()
		if err == io.EOF {
			break;
		}
		if err != nil {
			t.Fatal(err)
		}
		t.Log(currApp)
		if (!currApp.Hearthbeat) {
			i++
		}
		if i == 5 {
			return
		}
	}
	apps.CloseSend()

}