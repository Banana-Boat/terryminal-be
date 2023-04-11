package ws

import (
	"encoding/json"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type Message struct {
	Event string                 `json:"event" mapstructure:"event"`
	Data  map[string]interface{} `json:"data" mapstructure:"data"`
}

type LaunchClientData struct {
	ContainerName string `json:"containerName" mapstructure:"containerName"`
}
type LaunchServerData struct {
	Result bool `json:"result" mapstructure:"result"`
}

type CloseServerData struct {
	Result bool `json:"result" mapstructure:"result"`
}

type RunCmdClientData struct {
	Cmd string `json:"cmd" mapstructure:"cmd"`
}
type RunCmdServerData struct {
	Result string `json:"result" mapstructure:"result"`
}

func sendMessage(conn net.Conn, event string, data interface{}) {
	var _data map[string]interface{}
	mapstructure.Decode(data, &_data)

	msg, _ := json.Marshal(Message{
		Event: event,
		Data:  _data,
	})
	if err := wsutil.WriteServerMessage(conn, ws.OpText, msg); err != nil {
		log.Error().Err(err).Msg("failed to send message")
	}
}
