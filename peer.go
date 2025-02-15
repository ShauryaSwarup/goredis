package main

import (
	"bufio"
	"errors"
	"goredis/resp"
	"io"
	"log/slog"
	"net"
)

type Peer struct {
	conn       net.Conn
	cmdChan    chan Command
	respReader *resp.Resp
	respWriter *resp.Writer
}

func NewPeer(conn net.Conn, cmdChan chan Command) *Peer {
	return &Peer{
		conn:       conn,
		cmdChan:    cmdChan,
		respReader: resp.NewResp(bufio.NewReader(conn)),
		respWriter: resp.NewWriter(bufio.NewWriter(conn)),
	}
}

func (p *Peer) readLoop() error {
	defer p.conn.Close()

	// COMMAND DOCS
	val, err := p.respReader.ReadValue()
	if err != nil {
		return err
	}
	p.cmdChan <- Command{
		peer: p,
		Args: val,
	}

	for {
		value, err := p.respReader.ReadValue()
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Connection closed", "remoteAddr", p.conn.RemoteAddr())
				return nil
			}
			slog.Error("Read error", "Err:", err, "RemoteAddr", p.conn.RemoteAddr())
			return err
		}
		p.cmdChan <- Command{
			peer: p,
			Args: value,
		}
	}
}
