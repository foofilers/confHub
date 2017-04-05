package api

import (
	"github.com/coreos/etcd/embed"
	"github.com/Sirupsen/logrus"
	"os"
	"github.com/foofilers/confHub/conf"
	"github.com/foofilers/confHub/server"
	"gopkg.in/resty.v0"
	"testing"
)

var testcnf = `
version: 0.0.1b
jwtSecretKey: asdikasjdiowquaswdas9802uejbdsyu
pwdSecretKey: asdfhw21398712nw21ui873w121312kj
etcd:
  servers:
    - http://127.0.0.1:2379
  password: testRoot
`

var etcdTmpDir = "/tmp/confHubTest_etcd"
var RootPwd = "testRoot"
var ServerUrl = "http://127.0.0.1:9090"

func startServer() {
	var err error
	cfg := embed.NewConfig();
	cfg.Dir = etcdTmpDir
	_, err = embed.StartEtcd(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	conf.InitConf(testcnf)
	server.StartAsync("127.0.0.1:9090", true)
}

func TestMain(m *testing.M) {
	startServer()
	resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))
	m.Run()
	stopServer()
}

func stopServer() {
	logrus.Debug("remove etcd folder")
	os.RemoveAll(etcdTmpDir)
}

func checkHttpStatus(t *testing.T, resp *resty.Response, expected int) {
	if resp.StatusCode() != expected {
		t.Fatalf("status code should be %d but was %d", expected, resp.StatusCode())
	}
}
