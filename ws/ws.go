package ws

import (
	"net"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

type WebsocketConn struct {
	net.Conn

	closed int32
	wait   chan struct{}
}

func NewWebSocketConn(conn net.Conn) (c *WebsocketConn) {
	c = &WebsocketConn{
		Conn: conn,
		wait: make(chan struct{}),
	}
	return
}

func (p *WebsocketConn) Close() error {
	if atomic.SwapInt32(&p.closed, 1) == 1 {
		return nil
	}
	close(p.wait)
	return p.Conn.Close()
}

func (p *WebsocketConn) waitClose() {
	<-p.wait
}

// ConnectWebsocketServer :
// addr: ws://domain:port
func ConnectWebsocketServer(addr, origin string) (c net.Conn, err error) {
	cfg, err := websocket.NewConfig(addr, origin)
	if err != nil {
		return
	}
	cfg.Dialer = &net.Dialer{
		Timeout: time.Second * 10,
	}

	conn, err := websocket.DialConfig(cfg)
	if err != nil {
		return
	}
	c = NewWebSocketConn(conn)
	return
}
