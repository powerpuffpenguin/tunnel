package http3

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/powerpuffpenguin/tunnel/transport/http2"
	"github.com/quic-go/quic-go/http3"
)

type Client struct {
	client   http.Client
	rawURL   string
	from, to Addr
}

func NewClient(rawURL, from, to string, insecureSkipVerify bool) (c *Client, e error) {
	u, e := url.ParseRequestURI(to)
	if e != nil {
		return
	}
	if u.Scheme != `tcp` {
		e = errors.New(`not support to scheme: ` + u.Scheme)
		return
	}
	c = &Client{
		client: http.Client{
			Transport: &http3.RoundTripper{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecureSkipVerify,
				},
			},
		},
		rawURL: rawURL,
		from: Addr{
			network: `quic`,
			addr:    from,
		},
		to: Addr{
			network: u.Scheme,
			addr:    to,
		},
	}
	return
}
func (c *Client) Connect() (conn net.Conn, e error) {
	body, w := io.Pipe()
	resp, e := http2.Connect(&c.client,
		c.rawURL, c.to.addr,
		body,
	)
	if e != nil {
		w.Close()
		return
	}
	conn = http2.NewConn(resp.Body, w, c.from, c.to)
	return
}
