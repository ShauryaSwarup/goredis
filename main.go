package main

import (
	"fmt"
	"goredis/handler"
	"goredis/resp"
	"log"
	"log/slog"
	"net"
	"strings"
)

const defaultListenerAddr = ":6379"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	ln          net.Listener
	peers       map[*Peer]bool
	addPeerChan chan *Peer
	quitChan    chan struct{}
	cmdChan     chan Command
}

type Command struct {
	peer *Peer
	Args resp.Value
}

func NewServer(cfg Config) *Server {
	return &Server{
		Config:      cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		cmdChan:     make(chan Command, 1000),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	slog.Info("Server started", "ListenAddress: ", s.ListenAddress)

	// block starting (accept Connections)
	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Accept Error", "err:", err)
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) loop() {
	for {
		select {
		case cmd := <-s.cmdChan:
			if err := s.handleCmd(cmd); err != nil {
				slog.Error("Raw Message processing", "Error: ", err)
			}
		case <-s.quitChan:
			return
		case peer := <-s.addPeerChan:
			s.peers[peer] = true
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.cmdChan)
	s.addPeerChan <- peer
	slog.Info("New Peer connected", "remoteAddr: ", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		delete(s.peers, peer)
		fmt.Println("Peer at: ", peer.conn.RemoteAddr(), " Status: ", s.peers[peer])
		// slog.Error("Read peer error", "Error:", err, "RemoteAddr", conn.RemoteAddr())
	}
}

func (s *Server) handleCmd(cmd Command) error {
	slog.Info("HANDLING COMMAND",
		"type", cmd.Args.Typ,
		"args", cmd.Args,
	)

	// the request will always be decomposed into an array of Values
	// eg. SET name shaurya
	// This is nothing but an array of bulk strings.
	if cmd.Args.Typ != "array" {
		slog.Error("Incorrect type: ", "Type: ", cmd.Args.Typ)
	}
	if len(cmd.Args.Array) == 0 {
		slog.Error("Empty Request? ")
	}

	command := strings.ToUpper((cmd.Args.Array[0].Bulk))
	commandfunc, exists := handler.Handler[command]

	if !exists {
		slog.Error("Unsupported Func", "Func:", command)
		cmd.peer.respWriter.Write(resp.Value{Typ: "simplestring", Str: ""})
		return nil
	}

	res := commandfunc(cmd.Args.Array[1:])

	cmd.peer.respWriter.Write(res)

	return nil
}

func main() {
	cfg := Config{
		ListenAddress: ":6379",
	}
	server := NewServer(cfg)
	log.Fatal(server.Start())
}
