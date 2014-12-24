package service

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Connector interface {
	SendMessage(msg string) error
	SendData(data []byte) error
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	connectionChan chan *Conn
)

type Conn struct {
	ws       *websocket.Conn
	send     chan []byte
	recMsg   chan []byte
	recData  chan []byte
	doneChan chan bool
}

func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{
		ws:       ws,
		send:     make(chan []byte, 256),
		recMsg:   make(chan []byte, 256),
		recData:  make(chan []byte, 256),
		doneChan: make(chan bool),
	}
}

func (c *Conn) readPump() {
	defer func() {
		c.ws.Close()
		close(c.doneChan)
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		switch mt {
		case websocket.TextMessage:
			c.recMsg <- msg
		case websocket.BinaryMessage:
			c.recData <- msg
		}
	}
}

func (c *Conn) SendMessage(msg string) error {
	return c.write(websocket.TextMessage, []byte(msg))
}

func (c *Conn) SendData(data []byte) error {
	return c.write(websocket.BinaryMessage, data)
}

func (c *Conn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		close(c.doneChan)
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func init() {
	connectionChan = make(chan *Conn)
}
