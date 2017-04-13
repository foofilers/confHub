package ws

import (
	"gopkg.in/kataras/iris.v6/adaptors/websocket"
	"github.com/Sirupsen/logrus"
	"gopkg.in/kataras/iris.v6"
	"encoding/json"
)

func InitWs(app *iris.Framework) {
	ws := websocket.New(websocket.Config{
		ReadBufferSize:1024,
		WriteBufferSize:1024,
		Endpoint:"/ws",
	})
	ws.OnConnection(onConnect)
	app.Adapt(ws)
}

func onConnect(conn websocket.Connection) {
	logrus.Infof("New websocket connection id:%s", conn.ID())
	conn.OnMessage(func(msg []byte) {
		logrus.Infof("Received from %s:%s", conn.ID(), string(msg))
		wsCmd := &GenericWsCommand{}
		if err := json.Unmarshal(msg, wsCmd); err != nil {
			logrus.Error(err)
			conn.EmitError(err.Error())
			return
		}
		wsCmd.Exec(msg, conn)
	})
	conn.OnDisconnect(func() {
		logrus.Infof("%s Disconnected from websocket ", conn.ID())
		RemoveWatchConn(conn.ID())
	})
}