package ws

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
	Event string                 `json:"event" mapstructure:"event"`
	Data  map[string]interface{} `json:"data" mapstructure:"data"`
}

/* Message中Data具体定义 */
type StartClientData struct {
	PtyID string `json:"ptyID" mapstructure:"ptyID"`
}
type StartServerData struct {
	PtyID  string `json:"ptyID" mapstructure:"ptyID"`
	Result bool   `json:"result" mapstructure:"result"`
}

type EndClientData struct {
	PtyID string `json:"ptyID" mapstructure:"ptyID"`
}

type EndServerData struct {
	PtyID  string `json:"ptyID" mapstructure:"ptyID"`
	Result bool   `json:"result" mapstructure:"result"`
}

type RunCmdClientData struct {
	PtyID string `json:"ptyID" mapstructure:"ptyID"`
	Cmd   string `json:"cmd" mapstructure:"cmd"`
}
type RunCmdServerData struct {
	PtyID   string `json:"ptyID" mapstructure:"ptyID"`
	IsError bool   `json:"isError" mapstructure:"isError"`
	Result  string `json:"result" mapstructure:"result"`
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
