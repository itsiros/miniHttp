package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/tsironi93/miniHttp/internal/request"
	"github.com/tsironi93/miniHttp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

type HandlerFunc func(w io.Writer, req *request.Request) *response.HandleError

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) handleRequest(conn net.Conn) (*request.Request, error) {

	r, err := request.RequestFromReader(conn)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	if err := s.handleRequest(conn); err != nil {
		response.WriteError(conn, err)
		return
	}

	h := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn, response.SC200)
	response.WriteHeaders(conn, h)

}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			if s.closed.Load() {
				return
			}
		}
		go s.handle(conn)
	}
}

func Serve(port int, handler *HandlerFunc) (*Server, error) {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s := &Server{
		listener: ln,
	}

	go s.listen()

	return s, nil
}
