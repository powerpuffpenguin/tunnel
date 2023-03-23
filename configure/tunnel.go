package configure

type Tunnel struct {
	// Tunnel server URL
	URL string
	// Listen on local address
	From string
	// Forwarded server connection address
	To string
	// Do not verify tls certificate
	InsecureSkipVerify bool
	// Use quic to connect to the server URL
	Quic               bool
}

func LoadTunnel(filename string) (cnfs []Tunnel, e error) {
	var tmp []Tunnel
	e = loadObject(filename, &tmp)
	if e != nil {
		return
	}
	cnfs = tmp
	return
}
