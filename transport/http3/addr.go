package http3

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
