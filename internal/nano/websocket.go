package nano

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/cenkalti/log"
	"golang.org/x/net/websocket"
)

type Websocket struct {
	url  string
	conn *websocket.Conn
}

func NewWebsocket(wsURL string) *Websocket {
	return &Websocket{
		url: wsURL,
	}
}

func (w *Websocket) Connect() error {
	log.Debugf("connecting to websocket: %s", w.url)
	conn, err := websocket.Dial(w.url, "", "http://localhost/")
	if err != nil {
		return err
	}
	log.Debugf("connected to websocket: %s", w.url)
	w.conn = conn
	return nil
}

func (w *Websocket) Close() error {
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}

func (w *Websocket) Send(action, topic string, ack bool, options map[string]interface{}) error {
	m := map[string]interface{}{
		"action": action,
	}
	if topic != "" {
		m["topic"] = topic
	}
	if ack {
		m["ack"] = true
		m["id"] = strconv.Itoa(rand.Int()) // nolint: gosec
	}
	if len(options) > 0 {
		m["options"] = options
	}
	log.Debugf("sending websocket message: %#v", m)
	err := websocket.JSON.Send(w.conn, m)
	if err != nil {
		return nil
	}
	if ack {
		var ackMsg struct {
			Ack string `json:"ack"`
			ID  string `json:"id"`
		}
		err = w.Recv(&ackMsg)
		if err != nil {
			return err
		}
		if ackMsg.Ack != m["action"] || ackMsg.ID != m["id"] {
			return errInvalidAck
		}
	}
	return nil
}

var errInvalidAck = errors.New("invalid ack")

func (w *Websocket) Recv(msg interface{}) error {
	return websocket.JSON.Receive(w.conn, msg)
}
