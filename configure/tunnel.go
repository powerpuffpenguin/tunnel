package configure

type Tunnel struct {
	// tunnel server url
	URL string
	// listen local address
	From string
	// local to remote url
	To string
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
