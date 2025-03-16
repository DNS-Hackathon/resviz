package dig

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Query struct {
	Nameserver string `json:"Nameserver"`
	Transport  string `json:"Transport"`
	Qname      string `json:"Qname"`
	Qtype      string `json:"Qtype"`
	Port       string `json:"Port"`
	IpVersion  string `json:"IpVersion"`
	AA         bool   `json:"AA"`
	AD         bool   `json:"AD"`
	CD         bool   `json:"CD"`
	RD         bool   `json:"RD"`
	DO         bool   `json:"DO"`
	NoCrypto   bool   `json:"NoCrypto"`
	Nsid       bool   `json:"Nsid"`
	ShowQuery  bool   `json:"ShowQuery"`
	Reverse    bool   `json:"Reverse"`
	UDPsize    uint16 `json:"UDPsize"`
	Tsig       string `json:"Tsig"`
}

type WebQuery struct {
	Nameserver string `json:"Nameserver"`
	Transport  string `json:"Transport"`
	Qname      string `json:"Qname"`
	Qtype      string `json:"Qtype"`
	Port       string `json:"Port"`
	IpVersion  string `json:"IpVersion"`
	AA         string `json:"AA"`
	AD         string `json:"AD"`
	CD         string `json:"CD"`
	RD         string `json:"RD"`
	DO         string `json:"DO"`
	NoCrypto   string `json:"NoCrypto"`
	Nsid       string `json:"Nsid"`
	ShowQuery  string `json:"ShowQuery"`
	Reverse    string `json:"Reverse"`
	UDPsize    string `json:"UDPsize"`
	Tsig       string `json:"Tsig"`
}

type DigOut struct {
	Qname      string        `json:"Qname"`
	Query      *dns.Msg      `json:"Query"`    // message sent to name server
	Response   *dns.Msg      `json:"Response"` // response from nameserver
	RTT        time.Duration `json:"Round trip time"`
	Nameserver string        `json:"Nameserver"` // Name server IP
	QNSname    string        `json:"QNSname"`    // resolver name before translation
	ShowQuery  bool          `json:"ShowQuery"`
	MsgSize    int           `json:"Message Size"`
	Transport  string        `json:"Transport"`
}

// sanitize input data as precaution
func (q *Query) Sanitize() {
	q.Transport = strings.ToLower(q.Transport) // needs to be lower case.
	// if reverse is set and query name is a valid IP address, convert
	// qname to corresponding .arpa
	if q.Reverse {
		ip := net.ParseIP(q.Qname)

		if ip != nil {
			rev, err := dns.ReverseAddr(q.Qname)
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				q.Qname = rev
			}

		}
	}

}

// HTMX json-enc encodes checkbox values strings.
// If checked, defalts to "on" if value="true" attribute not set in HTML
func (wq *WebQuery) Parse() Query {
	var q Query
	q.Nameserver = wq.Nameserver
	q.Transport = wq.Transport
	q.Qname = wq.Qname
	q.Qtype = wq.Qtype
	q.Port = wq.Port
	q.IpVersion = wq.IpVersion
	q.AA = FixBool(wq.AA)
	q.AD = FixBool(wq.AD)
	q.CD = FixBool(wq.CD)
	q.RD = FixBool(wq.RD)
	q.DO = FixBool(wq.DO)
	q.NoCrypto = FixBool(wq.NoCrypto)
	q.Nsid = FixBool(wq.Nsid)
	q.ShowQuery = FixBool(wq.ShowQuery)
	q.Reverse = FixBool(wq.Reverse)
	udp, err := strconv.ParseUint(wq.UDPsize, 10, 32)
	if err != nil {
		q.UDPsize = 1232
	} else {
		if udp >= dns.MinMsgSize || udp <= dns.MaxMsgSize {
			q.UDPsize = uint16(udp)
		}
	}
	q.Tsig = wq.Tsig

	return q

}

func FixBool(s string) bool {
	// if the box is checked, it has a value, i.e. true. If not, false
	// "false added for inverse optons anv <input hidden....>"
	if s == "" || s == "false" {
		return false
	}
	return true

}

func GetSystemResolver(ipver string) string {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	// check for the first available server of right i version
	var ns string
	for _, ip := range conf.Servers {
		ns = net.ParseIP(ip).String()
		if ipver == "6" && strings.Contains(ns, ":") {
			break
		}
		if ipver == "4" && strings.Contains(ns, ".") {
			break
		}

	}
	return ns

}

// Harmonize lookup nameserver to always use IP:Port
// Check if valid IP. If not, assume hostname and look it up, selecting the
// first available ip of correct version
func (q *Query) GetLookupNS() string {
	var ns string

	// If no nameserver was passed, use system resolver
	if len(q.Nameserver) == 0 {
		q.Nameserver = GetSystemResolver(q.IpVersion)
	}

	ip := net.ParseIP(q.Nameserver)

	if ip != nil {
		if q.IpVersion == "6" {
			ns = "[" + q.Nameserver + "]:" + q.Port
		} else {
			ns = q.Nameserver + ":" + q.Port
		}
	} else {
		IPlist, err := net.LookupIP(q.Nameserver)
		if err != nil {
			log.Print(err)
		} else {
			for _, ip := range IPlist {
				if q.IpVersion == "6" {
					if strings.Count(ip.String(), ":") >= 2 {
						ns = "[" + ip.String() + "]:" + q.Port
						break
					}
				} else {
					// If address contains more than 1 ':', it's a V6 address. Go next.
					if strings.Count(ip.String(), ":") < 2 {
						ns = ip.String() + ":" + q.Port
						break
					}
				}
			}
		}

	}
	return ns
}
