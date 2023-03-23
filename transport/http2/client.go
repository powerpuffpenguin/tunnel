package http2

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

type Client struct {
	client   http.Client
	rawURL   string
	from, to Addr
}

func NewClient(rawURL,
	from, to string,
	insecureSkipVerify, h2c bool,
) (c *Client, e error) {
	var (
		transport http.RoundTripper
		addr      = Addr{
			addr: from,
		}
	)
	if h2c {
		addr.network = `h2c`
		transport = &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
			AllowHTTP: true,
		}
	} else {
		addr.network = `h2`
		transport = &http2.Transport{}
		if insecureSkipVerify {
			transport = &http2.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecureSkipVerify,
				},
			}
		}

	}

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
			Transport: transport,
		},
		rawURL: rawURL,
		from:   addr,
		to: Addr{
			network: u.Scheme,
			addr:    to,
		},
	}
	return
}

func (c *Client) Connect() (conn net.Conn, e error) {
	body, w := io.Pipe()
	resp, e := Connect(&c.client,
		c.rawURL, c.to.addr,
		body,
	)
	if e != nil {
		w.Close()
		return
	}
	conn = &clientConn{
		r:          resp.Body,
		w:          w,
		localAddr:  c.from,
		remoteAddr: c.to,
	}
	return
}

func Connect(client *http.Client,
	rawURL, to string,
	body io.Reader,
) (resp *http.Response, e error) {
	req, e := http.NewRequest(http.MethodPost, rawURL, body)
	if e != nil {
		return
	}
	req.Header.Set(`target`, base64.RawURLEncoding.EncodeToString([]byte(to)))
	resp, e = client.Do(req)
	if e != nil {
		return
	}
	if resp.StatusCode != http.StatusCreated {
		var b []byte
		if resp.Body != nil {
			b, _ = io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
		}
		if len(b) == 0 {
			e = errors.New(resp.Status)
		} else {
			e = errors.New(resp.Status + ` ` + string(b))
		}
		return
	} else if resp.Body == nil {
		e = errors.New(`body nil`)
		return
	}
	return
}

type Addr struct {
	network string
	addr    string
}

func (a Addr) Network() string {
	return a.network
}
func (a Addr) String() string {
	return a.addr
}

type clientConn struct {
	r                     io.ReadCloser
	w                     io.WriteCloser
	localAddr, remoteAddr net.Addr
}

func NewConn(r io.ReadCloser,
	w io.WriteCloser,
	localAddr, remoteAddr net.Addr,
) net.Conn {
	return &clientConn{
		r:          r,
		w:          w,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
	}
}
func (c *clientConn) Read(b []byte) (n int, err error) {
	n, err = c.r.Read(b)
	return
}
func (c *clientConn) Write(b []byte) (n int, err error) {
	n, err = c.w.Write(b)
	return
}

func (c *clientConn) Close() error {
	ew := c.w.Close()
	er := c.r.Close()
	if ew != nil {
		return ew
	}
	return er
}

func (c *clientConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *clientConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *clientConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *clientConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *clientConn) SetWriteDeadline(t time.Time) error {
	return nil
}
