package configure

type Server struct {
	// Listen address
	Listen string
	// TLS certificate, use h2c if not set
	CertFile string
	KeyFile  string
	// Router path
	Router string
	// quic or h2
	Quic bool
}

func LoadServers(filename string) (cnfs []Server, e error) {
	var tmp []Server
	e = loadObject(filename, &tmp)
	if e != nil {
		return
	}
	cnfs = tmp
	return
}
