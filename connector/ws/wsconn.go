package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amorwilliams/unigo/connector"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	Host string
	Port string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWsSocketProxy(host, port, config string) (socket Connector.Socket, err error) {
	Host = host
	Port = port
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	//TODO:解析配置

	socket = &WsSocket{
		pending: make(chan Connector.Connector),
	}

	return socket, nil
}

type WsSocket struct {
	pending chan Connector.Connector
}

func (ws *WsSocket) Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Upgrade:", err)
			return
		}
		defer func() {
			c.Close()
		}()

		wsc := NewConn(c)
		ws.pending <- wsc

		go wsc.writePump()
		wsc.readPump()
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", Port), nil)
	if err != nil {
		return err
	}
	return nil
}

func (ws *WsSocket) Stop() error {
	return nil
}

func (ws *WsSocket) Accept() chan Connector.Connector {
	return ws.pending
}

type State int32

const (
	ST_INITED State = iota
	ST_CLOSED
)

type WsConn struct {
	ws    *websocket.Conn
	state State

	send     chan []byte
	recive   chan []byte
	doneChan chan bool
}

func NewConn(ws *websocket.Conn) *WsConn {
	return &WsConn{
		ws:       ws,
		send:     make(chan []byte, 256),
		recive:   make(chan []byte, 256),
		doneChan: make(chan bool),
		state:    ST_INITED,
	}
}

func (c *WsConn) readPump() {
	defer func() {
		c.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		if mt == websocket.CloseMessage {
			c.Close()
		} else if mt == websocket.TextMessage {
			c.recive <- []byte(msg)
		} else if mt == websocket.BinaryMessage {
			c.recive <- msg
		}
	}
}

func (c *WsConn) Recive() chan []byte {
	return c.recive
}

func (c *WsConn) Send(data []byte) error {
	c.send <- data
	return nil
}

func (c *WsConn) Close() {
	if c.state == ST_CLOSED {
		return
	}

	c.state = ST_CLOSED
	close(c.doneChan)
	c.ws.Close()
}

func (c *WsConn) Done() chan bool {
	return c.doneChan
}

func (c *WsConn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *WsConn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.BinaryMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-c.doneChan:
			return
		}
	}
}
