package configure

type Server struct {
	Listen   string
	CertFile string
	KeyFile  string
	Router   string
	
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
