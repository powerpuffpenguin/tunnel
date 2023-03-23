package main

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/powerpuffpenguin/tunnel/configure"
	"github.com/powerpuffpenguin/tunnel/transport/http2"
	"github.com/powerpuffpenguin/tunnel/transport/http3"
	"github.com/spf13/cobra"
)

func server() *cobra.Command {
	var (
		cnf string
	)
	cmd := &cobra.Command{
		Use:   "server",
		Short: "run tunnel server",
		Run: func(cmd *cobra.Command, args []string) {
			cnfs, e := configure.LoadServers(cnf)
			if e != nil {
				log.Fatalln(e)
			}
			srvs := make([]Server, 0, len(cnfs))
			for i := 0; i < len(cnfs); i++ {
				s, e := NewServer(&cnfs[i])
				if e != nil {
					log.Fatalln(e)
				}
				srvs = append(srvs, s)
			}
			var wg sync.WaitGroup
			wg.Add(len(srvs))
			for _, s := range srvs {
				go func(s Server) {
					e := s.Serve()
					if e != nil {
						log.Fatalln(e)
					}
					wg.Done()
				}(s)
			}
			wg.Wait()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&cnf, `cnf`, `c`, filepath.Join(BasePath(), `etc`, `server.jsonnet`), `configure file`)
	return cmd
}

type Server interface {
	Serve() error
	Close() error
}

func NewServer(cnf *configure.Server) (s Server, e error) {
	if cnf.CertFile == `` || cnf.KeyFile == `` {
		s, e = http2.New(cnf.Listen, cnf.CertFile, cnf.KeyFile, cnf.Router)
		if e != nil {
			return
		}
		log.Println(`h2c listen`, cnf.Listen, cnf.Router)
	} else if cnf.Quic {
		s, e = http3.New(cnf.Listen, cnf.CertFile, cnf.KeyFile, cnf.Router)
		if e != nil {
			return
		}
		log.Println(`quic listen`, cnf.Listen, cnf.Router)
	} else {
		s, e = http2.New(cnf.Listen, cnf.CertFile, cnf.KeyFile, cnf.Router)
		if e != nil {
			return
		}
		log.Println(`h2 listen`, cnf.Listen, cnf.Router)
	}
	return
}
