package tuya

import (
	"encoding/hex"
	"log"
	"net"
)

type Server struct {
	listener net.PacketConn
	debug    bool
}

func NewListener(debug bool) *Server {
	l, err := net.ListenPacket("udp", "0.0.0.0:6666")
	if err != nil {
		return nil
	}

	if debug {
		log.Println("listen on address", l.LocalAddr())
	}

	server := &Server{
		listener: l,
		debug:    debug,
	}

	return server
}

func (s *Server) Close() error {
	log.Println("close listener")
	return s.listener.Close()
}

func (s *Server) Receiver(e *Engine) {
	for {
		buffer := make([]byte, 0x1024)
		n, addr, err := s.listener.ReadFrom(buffer)
		if err != nil {
			continue
		}

		log.Println("receive packet from", addr)

		log.Println("packet lenth", n)
		if n > 16 {
			hex.Dump(buffer)
		}
	}
}
