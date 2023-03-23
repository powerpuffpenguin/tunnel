package main

import (
	"log"
	"net"
	"path/filepath"
	"sync"
	"time"

	"github.com/powerpuffpenguin/tunnel/configure"
	"github.com/powerpuffpenguin/tunnel/pool"
	"github.com/powerpuffpenguin/tunnel/transport/http2"
	"github.com/spf13/cobra"
)

func tunnel() *cobra.Command {
	var (
		cnf string
	)
	cmd := &cobra.Command{
		Use:   "tunnel",
		Short: "run tunnel",
		Run: func(cmd *cobra.Command, args []string) {
			cnfs, e := configure.LoadTunnel(cnf)
			if e != nil {
				log.Fatalln(e)
			}
			srvs := make([]*Tunnel, 0, len(cnfs))
			for i := 0; i < len(cnfs); i++ {
				s, e := NewTunnel(&cnfs[i])
				if e != nil {
					log.Fatalln(e)
				}
				srvs = append(srvs, s)
			}
			var wg sync.WaitGroup
			wg.Add(len(srvs))
			for _, s := range srvs {
				go func(s *Tunnel) {
					s.Serve()
					wg.Done()
				}(s)
			}
			wg.Wait()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&cnf, `cnf`, `c`, filepath.Join(BasePath(), `etc`, `tunnel.jsonnet`), `configure file`)
	return cmd
}

type Client interface {
	Connect() (net.Conn, error)
}
type Tunnel struct {
	l      net.Listener
	client Client
}

func NewTunnel(cnf *configure.Tunnel) (t *Tunnel, e error) {
	c, e := http2.NewClient(cnf.URL, cnf.From, cnf.To)
	if e != nil {
		return
	}

	l, e := net.Listen(`tcp`, cnf.From)
	if e != nil {
		return
	}
	log.Println(`tunnel`, cnf.URL, cnf.From, `->`, cnf.To)

	t = &Tunnel{
		l:      l,
		client: c,
	}
	return
}
func (t *Tunnel) Serve() error {
	var tempDelay time.Duration
	for {
		rw, err := t.l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("http: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		tempDelay = 0
		go t.serve(rw)
	}
}
func (t *Tunnel) serve(c0 net.Conn) {
	c1, e := t.client.Connect()
	if e != nil {
		log.Println(e)
		c0.Close()
		return
	}

	done := make(chan int, 1)
	go t.copy(c0, c1, done)
	go t.copy(c1, c0, done)

	<-done
	time.Sleep(time.Second)
	c0.Close()
	c1.Close()
}
func (t *Tunnel) copy(dst, src net.Conn, done chan int) {
	var (
		b      = pool.GetBytes()
		n      int
		er, ew error
	)
	for {
		n, er = src.Read(b)
		if n > 0 {
			_, ew = dst.Write(b[:n])
			if ew != nil {
				break
			}
		}
		if er != nil {
			break
		}
	}
	pool.PutBytes(b)
	done <- 1
}
