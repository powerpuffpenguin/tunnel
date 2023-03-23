package http2

import (
	"encoding/base64"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/powerpuffpenguin/tunnel/pool"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	server            http.Server
	l                 net.Listener
	h2c               bool
	certFile, keyFile string
}

func New(listen, certFile, keyFile, router string) (s *Server, e error) {
	l, e := net.Listen(`tcp`, listen)
	if e != nil {
		return
	}
	s = &Server{
		l:        l,
		certFile: certFile,
		keyFile:  keyFile,
		h2c:      certFile == `` || keyFile == ``,
	}
	var (
		http2Server http2.Server
		mux         = http.NewServeMux()
	)
	mux.HandleFunc(router, routerF)
	if s.h2c {
		s.server.Handler = h2c.NewHandler(mux, &http2Server)
	} else {
		s.server.Handler = mux
	}
	e = http2.ConfigureServer(&s.server, &http2Server)
	if e != nil {
		l.Close()
		return
	}

	return
}
func (s *Server) Serve() (e error) {
	if s.h2c {
		e = s.server.Serve(s.l)
	} else {
		e = s.server.ServeTLS(s.l, s.certFile, s.keyFile)
	}
	return
}
func (s *Server) TLS() bool {
	return !s.h2c
}
func routerF(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor < 2 || r.Method != http.MethodPost || r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	target := r.Header.Get(`target`)
	b, e := base64.RawURLEncoding.DecodeString(target)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set(`Content-Type`, `text/plain; charset=utf-8`)
		w.Write([]byte(e.Error()))
		return
	}
	u, e := url.ParseRequestURI(string(b))
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set(`Content-Type`, `text/plain; charset=utf-8`)
		w.Write([]byte(e.Error()))
		return
	} else if u.Scheme != `tcp` {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set(`Content-Type`, `text/plain; charset=utf-8`)
		w.Write([]byte(`not support to scheme: ` + u.Scheme))
		return
	}

	ctx := r.Context()
	var dialer net.Dialer
	c, e := dialer.DialContext(ctx, `tcp`, u.Host)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set(`Content-Type`, `text/plain; charset=utf-8`)
		w.Write([]byte(e.Error()))
		return
	}
	defer c.Close()
	w.WriteHeader(http.StatusCreated)
	f := w.(http.Flusher)
	f.Flush()

	done := make(chan int, 1)
	go func() {
		var (
			b      = pool.GetBytes()
			n      int
			er, ew error
		)
		for {
			n, er = c.Read(b)
			if n > 0 {
				_, ew = w.Write(b[:n])
				if ew != nil {
					break
				}
				f.Flush()
			}
			if er != nil {
				break
			}
		}
		pool.PutBytes(b)
		done <- 1
	}()

	go func() {
		var (
			b      = pool.GetBytes()
			n      int
			er, ew error
		)
		for {
			n, er = r.Body.Read(b)
			if n > 0 {
				_, ew = c.Write(b[:n])
				if ew != nil {
					break
				}
			}
			if er != nil {
				break
			}
		}
		pool.PutBytes(b)
		done <- 2
	}()
	// wait any error
	<-done

	// waiting to read unread data in the network
	time.Sleep(time.Second)
}
