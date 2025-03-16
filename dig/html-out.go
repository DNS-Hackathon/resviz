package dig

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const (
	il    string = "info" // style selector to use with [class^="info-"] in css for styling info links
	hxget string = "dig/info/"
)

var copyicon = LoadSVG("html/images/copyIcon.svg")
var expandicon = LoadSVG("html/images/expandIcon.svg")
var version = "0.5"

func (r *DigOut) ToHTML() string {

	var out, banner, header, question, answer, authority, opt, additional, footer string
	var qheader, qopt, qquestion, qanswer, qauthority, qadditional string

	banner = "<tr>\n"
	banner += "<td colspan='5'>; <<>> NerDiG " + version + " <<>></td>\n"
	banner += "</tr>\n"

	qheader += headerToHTML(r.Query)
	qquestion = questonToHTML(r.Query)
	qanswer = answerToHTML(r.Query)
	qauthority = authorityToHTML(r.Query)
	qopt += optToHTML(r.Query)
	qadditional += additionalToHTML(r.Query)

	header += headerToHTML(r.Response)
	question = questonToHTML(r.Response)
	answer = answerToHTML(r.Response)
	authority = authorityToHTML(r.Response)
	opt = optToHTML(r.Response)
	additional = additionalToHTML(r.Response)
	footer = footerToHTML(r)

	// add the expand icon
	out += "<div class='iconBox'>\n"
	out += "<label for='wide-term'>\n"
	out += "<span>" + expandicon + "</span>\n"
	out += "</label>\n"
	out += "</div>\n"
	out += "<table class='fade-in'>\n"
	out += banner

	if r.ShowQuery {
		out += qheader
		out += qopt
		out += qquestion
		out += qanswer
		out += qauthority
		out += qadditional
	}
	out += header
	if r.Response.Truncated {
		out += "<tr>\n<td colspan='5'>\n;; WARNING: Response truncated (TC). Retry query over TCP</td>\n</tr>\n"
	}
	if r.Query.RecursionDesired && !r.Response.RecursionAvailable {
		out += "<tr>\n<td colspan='5'>\n;; WARNING: recursion requested but not available</td>\n</tr>\n"
	}
	var lf string
	if len(r.Response.Answer) > 0 {
		lf = dns.TypeToString[r.Response.Answer[len(r.Response.Answer)-1].Header().Rrtype]
	}
	if dns.TypeToString[r.Query.Question[0].Qtype][1:] == "XFR" && lf != "SOA" {
		out += "<tr>\n<td colspan='5' class='attention'>\n;; WARNING: AXFR/IXFR not fully supported. Answer was truncated!</td>\n</tr>\n"
	}
	// Manual fix for spacer here. Niceify later...
	out += "<tr>\n<td colspan='5' class='spacer'></td></tr>\n"
	out += opt
	out += question
	out += answer
	out += authority
	out += additional
	out += footer
	out += "</table>\n"

	return out
}

func headerToHTML(msg *dns.Msg) string {

	opcode := dns.OpcodeToString[msg.Opcode]
	rcode := dns.RcodeToString[msg.Rcode]
	id := strconv.Itoa(int(msg.Id))

	setflag := map[string]string{
		"QR": "unset",
		"RD": "unset",
		"RA": "unset",
		"AA": "unset",
		"AD": "unset",
		"CD": "unset",
		"TC": "unset",
	}

	// Sending vs receiving
	direction := "Sending:"
	if msg.Response {
		setflag["QR"] = "set"
		direction = "Got answer:"
	}

	if msg.RecursionDesired {
		setflag["RD"] = "set"
	}

	if msg.RecursionAvailable {
		setflag["RA"] = "set"
	}

	if msg.Authoritative {
		setflag["AA"] = "set"
	}

	if msg.AuthenticatedData {
		setflag["AD"] = "set"
	}

	if msg.CheckingDisabled {
		setflag["CD"] = "set"
	}

	if msg.Truncated {
		setflag["TC"] = "set"
	}

	flags := []string{"QR", "RD", "RA", "AA", "AD", "CD", "TC"}

	var header string

	header += "<tr>\n"
	header += "<td colspan='5'>;; " + direction + "</td>\n"
	header += "</tr>\n"
	header += "<tr>\n"
	header += "<td colspan='5'>\n"
	header += ";; ->>HEADER<<-" + hxwrap(";; opcode: "+opcode+",", "span", "opcode", []string{opcode, il})
	header += hxwrap("status: "+rcode+",", "span", "rcode", []string{rcode, il})
	header += hxwrap("id: "+id, "span", "message-id", []string{il})
	header += "</td>\n"
	header += "</tr>\n"
	header += "<tr>\n"
	header += "<td colspan='5'>\n"
	header += hxwrap(";; flags: ", "span", "flags", []string{"flags", il})
	for _, flag := range flags {
		header += hxwrap(flag, "span", "flag-"+flag, []string{flag, il, setflag[flag]})
	}
	header += "<span class='sepcolon'>;</span>"
	header += hxwrap("QUERY: "+strconv.Itoa(len(msg.Question))+",", "span", "count-QUERY", []string{il})
	header += hxwrap("ANSWER: "+strconv.Itoa(len(msg.Answer))+",", "span", "count-ANSWER", []string{il})
	header += hxwrap("AUTHORITY: "+strconv.Itoa(len(msg.Ns))+",", "span", "count-AUTHORITY", []string{il})
	header += hxwrap("ADDITIONAL: "+strconv.Itoa(len(msg.Extra)), "span", "count-ADDITIONAL", []string{il})
	header += "</td>\n"
	header += "</tr>\n"
	//header += "<tr>\n<td colspan='5' class='spacer'></td></tr>\n"

	return header
}

func questonToHTML(msg *dns.Msg) string {
	var question string

	var status string
	if len(msg.Question) < 1 {
		status = "unset"
	}

	question += "<tr class='" + status + "'>\n"
	question += "<td colspan='5'>\n"
	question += hxwrap(";; QUESTION SECTION:", "span", "section-QUESTION", []string{il})
	question += "</td>\n"
	question += "</tr>\n"
	question += "<td colspan='5'>\n"
	for _, q := range msg.Question {
		question += hxwrap(strings.TrimSpace(q.String()), "span", "placeholder", []string{il})
	}
	question += "</td>\n"
	question += "</tr>\n"
	question += "<tr>\n<td colspan='5' class='spacer'></td></tr>\n"

	return question
}

func answerToHTML(msg *dns.Msg) string {
	var answer string

	var status string
	if len(msg.Answer) < 1 {
		status = "unset"
	}

	answer += "<tr class='" + status + "'>\n"
	answer += "<td colspan='5'>\n"
	answer += hxwrap(";; ANSWER SECTION:", "span", "ANSWER-section", []string{il})
	answer += "</td>\n"
	answer += "</tr>\n"

	for _, a := range msg.Answer {

		head := *a.Header()

		answer += "<tr>\n"
		answer += hxwrap(head.Name, "td", "RR-owner-name", []string{il})
		answer += hxwrap(strconv.FormatUint(uint64(head.Ttl), 10), "td", "RR-ttl", []string{il})
		answer += hxwrap(dns.ClassToString[head.Class], "td", "RR-class", []string{il})
		answer += hxwrap(dns.Type(head.Rrtype).String(), "td", "RR-"+dns.Type(head.Rrtype).String(), []string{il})

		answer += rdatawrap(a, "RR-"+dns.Type(head.Rrtype).String())
		answer += "</tr>\n"
		/*
			rdata := ""
			for i := 1; i <= dns.NumField(a); i++ {
				rdata += dns.Field(a, i) + " "
				//out += dns.Field(a, i) + " "
				//out += "<div class=\"" + rfields[i-1] + "\">" + dns.Field(a, i) + "</div>\n"
				//fmt.Printf("field %v - %#v ", i, dns.Field(a, i))
			}
			answer += hxwrap(rdata, "td", "rdata-"+dns.Type(head.Rrtype).String(), []string{il, "rdata"})
		*/
	}
	answer += "<tr>\n<td colspan='5' class='spacer'></td></tr>\n"

	return answer
}

func authorityToHTML(msg *dns.Msg) string {

	var authority string

	var status string
	if len(msg.Ns) < 1 {
		status = "unset"
	}

	authority += "<tr class='" + status + "'>\n"
	authority += "<td colspan='5'>\n"
	authority += hxwrap(";; AUTHORITY SECTION:", "span", "AUTHORITY-section", []string{il})
	authority += "</td>\n"
	authority += "</tr>\n"
	for _, a := range msg.Ns {

		head := *a.Header()

		authority += "<tr>\n"
		authority += hxwrap(head.Name, "td", "owner-name", []string{il})
		authority += hxwrap(strconv.FormatUint(uint64(head.Ttl), 10), "td", "ttl", []string{il})
		authority += hxwrap(dns.ClassToString[head.Class], "td", "class", []string{il})
		authority += hxwrap(dns.Type(head.Rrtype).String(), "td", "RR-"+dns.Type(head.Rrtype).String(), []string{il})

		authority += rdatawrap(a, dns.Type(head.Rrtype).String())
		authority += "</tr>\n"
	}
	authority += "<tr>\n<td colspan='5' class='spacer'></td></tr>\n"

	return authority
}

func additionalToHTML(msg *dns.Msg) string {
	var additional string

	var status string
	if len(msg.Extra) < 2 {
		status = "unset"
	}

	additional += "<tr class='" + status + "'>\n"
	additional += "<td colspan='5'>\n"
	additional += hxwrap(";; ADDITIONAL SECTION:", "span", "ADDITIONAL-section", []string{il})
	additional += "</td>\n"
	additional += "</tr>\n"

	for _, e := range msg.Extra {
		head := *e.Header()

		if dns.Type(head.Rrtype).String() != "OPT" {

			additional += "<tr>\n"
			additional += hxwrap(head.Name, "td", "owner-name", []string{il})
			additional += hxwrap(strconv.FormatUint(uint64(head.Ttl), 10), "td", "ttl", []string{il})
			additional += hxwrap(dns.ClassToString[head.Class], "td", "class", []string{il})
			additional += hxwrap(dns.Type(head.Rrtype).String(), "td", "RR-"+dns.Type(head.Rrtype).String(), []string{il})

			additional += rdatawrap(e, dns.Type(head.Rrtype).String())
			additional += "</tr>\n"

		}
	}
	additional += "<tr>\n<td colspan='5' class='spacer'></td></tr>\n"

	return additional
}

func optToHTML(msg *dns.Msg) string {
	var opt string

	var status string
	if len(msg.Extra) < 1 {
		status = "unset"
	}

	opt += "<tr class='" + status + "'>\n"
	opt += "<td colspan='5'>\n"
	opt += hxwrap(";; OPT PSEUDOSECTION:", "span", "OPT-section", []string{il})
	opt += "</td>"
	opt += "</tr>\n"

	for _, e := range msg.Extra {
		head := *e.Header()

		if dns.Type(head.Rrtype).String() == "OPT" {

			f := e.(*dns.OPT)

			//
			// Pretty format the OPT section like
			//
			opt += "<tr>\n"
			opt += "<td colspan='5'>\n"
			opt += hxwrap("; EDNS: version "+strconv.Itoa(int(f.Version()))+"; ", "span", "OPT-EDNS", []string{il})

			var fs string
			if f.Do() {
				fs = "flags: do; "
			} else {
				fs = "flags:; "
			}
			opt += hxwrap(fs, "span", "OPT-doflag", []string{il})

			if f.Hdr.Ttl&0x7FFF != 0 {
				ms := fmt.Sprintf("MBZ: 0x%04x, ", f.Hdr.Ttl&0x7FFF)
				opt += hxwrap(ms, "span", "MBZ", []string{il})
			}

			opt += hxwrap("udp: "+strconv.Itoa(int(f.UDPSize())), "span", "OPT-udp", []string{il})
			opt += "</td>\n"
			opt += "</tr>\n"

			for _, o := range f.Option {
				var s string
				opt += "<tr>\n"
				opt += "<td colspan='5'>\n"
				switch o.(type) {
				case *dns.EDNS0_NSID:
					to := o.(*dns.EDNS0_NSID)
					s += "\n; NSID: " + to.String()
					h, e := PackNSID(to)
					var r string
					if e == nil {
						for _, c := range h {
							r += "(" + string(c) + ")"
						}
						s += "  " + r
					}
				case *dns.EDNS0_SUBNET:
					s += "; SUBNET: " + o.String()
				case *dns.EDNS0_COOKIE:
					s += "; COOKIE: " + o.String()
				case *dns.EDNS0_EXPIRE:
					s += "; EXPIRE: " + o.String()
				case *dns.EDNS0_TCP_KEEPALIVE:
					s += "; KEEPALIVE: " + o.String()
				case *dns.EDNS0_UL:
					s += "; UPDATE LEASE: " + o.String()
				case *dns.EDNS0_LLQ:
					s += "; LONG LIVED QUERIES: " + o.String()
				case *dns.EDNS0_DAU:
					s += "; DNSSEC ALGORITHM UNDERSTOOD: " + o.String()
				case *dns.EDNS0_DHU:
					s += "; DS HASH UNDERSTOOD: " + o.String()
				case *dns.EDNS0_N3U:
					s += "; NSEC3 HASH UNDERSTOOD: " + o.String()
				case *dns.EDNS0_LOCAL:
					s += "; LOCAL OPT: " + o.String()
				case *dns.EDNS0_PADDING:
					s += "; PADDING: " + o.String()
				case *dns.EDNS0_EDE:
					s += "; EDE: " + o.String()
				case *dns.EDNS0_ESU:
					s += "; ESU: " + o.String()
				}
				opt += hxwrap(s, "span", "EDNS-placeholder", []string{il})
				opt += "</td>\n"
				opt += "</tr>\n"
			}

		}
	}
	opt += "<tr>\n<td colspan='5' class='spacer'></td></tr>\n"

	return opt
}

func footerToHTML(r *DigOut) string {
	footer := "<tr>\n"
	footer += "<td colspan='5'>\n"
	// divide the Nanoseconds by 1e6 to get the Milliseconds as a int64
	footer += ";; " + hxwrap("Query time: "+strconv.Itoa(int(r.RTT)/1e6)+" ms", "span", "query-time", []string{il})
	footer += "</td>\n"
	footer += "</tr>\n"
	footer += "<tr>\n"
	footer += "<td colspan='5'>\n"
	footer += ";; " + hxwrap("SERVER: "+r.Nameserver+"("+r.QNSname+") ("+r.Transport[:3]+")", "span", "query-server", []string{il})
	footer += "</td>\n"
	footer += "</tr>\n"
	footer += "<tr>\n"
	footer += "<td colspan='5'>\n"
	footer += ";; WHEN: " + time.Now().Format(time.UnixDate)
	footer += "</td>\n"
	footer += "</tr>\n"
	footer += "<tr>\n"
	footer += "<td colspan='5'>\n"
	footer += ";; " + hxwrap("MSG SIZE: "+strconv.Itoa(r.Response.Len()), "span", "msg-size", []string{il})
	footer += "</td>\n"
	footer += "</tr>\n"

	return footer

}

func hxwrap(txt, tag, rfield string, cs []string) string {
	// Add classes to rr parts for styling and htmx links for info retrieval
	class := strings.Join(cs, " ")
	ws := "<" + tag + " class='" + class + "' "
	ws += "hx-get='" + hxget + rfield + "' "
	ws += "hx-target='#infobox' "
	ws += "hx-swap='innerHTML' "
	ws += ">\n" + txt + "\n</" + tag + ">\n"

	return ws

}

func rdatawrap(rr dns.RR, rtype string) string {

	rdata := ""
	wrr := ""

	// Add info for rdata
	switch rtype {
	case "TXT":
		for i := 1; i <= dns.NumField(rr); i++ {
			rdata += "&quot;" + dns.Field(rr, i) + "&quot; "
		}
		wrr += hxwrap(rdata, "td", "rdata-"+rtype, []string{il, "rdata"})
	default:
		for i := 1; i <= dns.NumField(rr); i++ {
			rdata += dns.Field(rr, i) + " "
		}
		wrr += hxwrap(rdata, "td", rtype, []string{il, "rdata"})
	}
	return wrr
}

func (q *Query) ToCLI() string {

	qs := "<div id='digcli' class='digcli-box' hx-swap-oob='outerHTML'>\n"
	qs += "<span id='cpcmd' onclick='copycmd(\"cpcmd\")'>"

	qs += "dig "

	if q.Nameserver != "" {
		qs += "@" + q.Nameserver + " "
	}

	if q.Reverse {
		qs += "-x "
	}

	qs += q.Qname + " "
	qs += q.Qtype + " "

	if q.Port != "53" {
		qs += "-p " + q.Port + " "
	}

	if q.IpVersion == "6" {
		qs += "-6 "

	}

	if q.AA {
		qs += "+aa "
	}

	if !q.AD {
		qs += "+noad "
	}

	if !q.RD {
		qs += "+nord "
	}

	if q.CD {
		qs += "+cd "
	}

	if q.ShowQuery {
		qs += "+qr "
	}

	if q.DO {
		qs += "+dnssec "
	}

	if q.NoCrypto {
		qs += "+nocrypto "
	}

	if q.Nsid {
		qs += "+nsid "
	}

	// from 9.18 default UDP buffer size is 1232
	// https://kb.isc.org/docs/behavior-dig-versions-edns-bufsize
	if q.UDPsize != 1232 {
		qs += "+bufsize=" + strconv.FormatUint(uint64(q.UDPsize), 10) + " "
	}

	if q.Transport == "tcp" {
		qs += "+tcp "
	}

	if q.Tsig != "" {
		qs += "-y " + q.Tsig + ""
	}

	// Add the copy icon
	//qs += LoadSVG("html/images/copyIcon.svg")
	qs += "<div class='iconBox'>" + copyicon + "</span>\n"

	qs += "</span>"
	qs += "</div>\n"

	return qs
}

func PackNSID(e *dns.EDNS0_NSID) ([]byte, error) {
	h, err := hex.DecodeString(e.Nsid)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func LoadSVG(file string) string {
	img, err := os.ReadFile(file)
	if err != nil {
	}
	return string(img)
}
