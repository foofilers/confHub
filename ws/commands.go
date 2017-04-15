package ws

import (
	"gopkg.in/kataras/iris.v6/adaptors/websocket"
	"encoding/json"
	"github.com/Sirupsen/logrus"
)

type Operation int

const (
	WATCH = iota
)

type WsCommand interface {
	Exec([]byte, websocket.Connection)
}

type WatchCommand struct {
	GenericWsCommand
	Data struct {
		     Applications []string `json:"applications"`
	     }
}

type GenericWsCommand struct {
	Op Operation `json:"operation"`
}

func (cmd *GenericWsCommand) Exec(data []byte, conn websocket.Connection) {
	var command WsCommand
	switch cmd.Op {
	case WATCH:
		command = &WatchCommand{}
	default:
		conn.EmitError("No Operation found")
		return
	}

	if err := json.Unmarshal(data, command); err != nil {
		conn.EmitError(err.Error())
		return
	}
	command.Exec(data, conn)
}

func (cmd *WatchCommand) Exec(data []byte, conn websocket.Connection) {
	logrus.Debugf("WatchCommand Exec %+v", cmd)
	for _, app := range cmd.Data.Applications {
		AddWatchConnAndStartNotifier(app, conn)
	}
}
