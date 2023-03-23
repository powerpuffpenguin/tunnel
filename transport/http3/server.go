package http3

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/powerpuffpenguin/tunnel/transport/http2"
	"github.com/quic-go/quic-go/http3"
)

type Server struct {
	quicServer *http3.Server

	udpConn *net.UDPConn
}

func New(listen, certFile, keyFile, router string) (s *Server, e error) {
	certs := make([]tls.Certificate, 1)
	certs[0], e = tls.LoadX509KeyPair(certFile, keyFile)
	if e != nil {
		return
	}
	if listen == "" {
		listen = ":https"
	}
	// Open the listeners
	udpAddr, e := net.ResolveUDPAddr("udp", listen)
	if e != nil {
		return
	}
	udpConn, e := net.ListenUDP("udp", udpAddr)
	if e != nil {
		return
	}
	// Start the servers
	mux := http.NewServeMux()
	mux.HandleFunc(router, http2.Router)
	s = &Server{
		quicServer: &http3.Server{
			TLSConfig: &tls.Config{
				Certificates: certs,
			},
			Handler: mux,
		},
		udpConn: udpConn,
	}
	return
}
func (s *Server) Serve() (e error) {
	e = s.quicServer.Serve(s.udpConn)
	return
}

func (s *Server) Close() error {
	e0 := s.quicServer.Close()
	e1 := s.udpConn.Close()
	if e0 != nil {
		return e0
	}
	return e1
}
