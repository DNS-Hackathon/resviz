package dig

var DefaultDNSquery = Query{
	Nameserver: "",
	Transport:  "udp",
	Qname:      "",
	Qtype:      "",
	Port:       "53",
	IpVersion:  "4",
	AA:         false,
	AD:         false,
	CD:         false,
	RD:         false,
	DO:         false,
	NoCrypto:   false,
	Nsid:       false,
	ShowQuery:  false,
	UDPsize:    1232,
	Tsig:       "",
}

func (q *Query) QminWalk(child, parent string) {
	q.Nameserver = parent
	q.Qname = child
	q.Qtype = "soa"
}
