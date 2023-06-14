package api

import (
	"encoding/json"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

// websocket message定义
type Message struct {
	PtyId string                 `json:"ptyId" mapstructure:"ptyId"`
	Event string                 `json:"event" mapstructure:"event"`
	Data  map[string]interface{} `json:"data" mapstructure:"data"`
}

/* Message中Data具体定义 */
type StartServerData struct {
	Result bool `json:"result" mapstructure:"result"`
}

type EndServerData struct {
	Result bool `json:"result" mapstructure:"result"`
}

type RunCmdClientData struct {
	Cmd string `json:"cmd" mapstructure:"cmd"`
}
type RunCmdServerData struct {
	IsError bool   `json:"isError" mapstructure:"isError"`
	Result  string `json:"result" mapstructure:"result"`
}

func sendMessage(conn net.Conn, ptyId string, event string, data interface{}) {
	var _data map[string]interface{}
	mapstructure.Decode(data, &_data)

	msg, _ := json.Marshal(Message{
		PtyId: ptyId,
		Event: event,
		Data:  _data,
	})
	if err := wsutil.WriteServerMessage(conn, ws.OpText, msg); err != nil {
		log.Error().Err(err).Msg("failed to send message")
	}
}
