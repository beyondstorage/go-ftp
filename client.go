package ftp

import (
	"crypto/tls"
	"net"
	"net/textproto"
)

type Client struct {
	conn *textproto.Conn
}

type Config struct {
	TLSConfig *tls.Config
}

func NewClient(addr string, cfg Config) (client *Client, err error) {
	var conn net.Conn
	if cfg.TLSConfig == nil {
		conn, err = net.Dial("tcp", addr)
	} else {
		conn, err = tls.Dial("tcp", addr, cfg.TLSConfig)
	}
	if err != nil {
		return nil, err
	}

	client = &Client{
		conn: textproto.NewConn(conn),
	}

	_, _, err = client.conn.ReadResponse(StatusReady)
	if err != nil {
		_ = client.Quit()
		return nil, err
	}
	return
}

func (c *Client) Quit() (err error) {
	defer c.conn.Close()

	_, err = c.conn.Cmd("QUIT")
	if err != nil {
		return err
	}
	return nil
}
