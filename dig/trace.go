package dig

import "fmt"

func Qmin(domain string) DigOut {

	var out DigOut
	// query ROOT-servers for TLD
	query := DefaultDNSquery
	query.QminWalk("192.36.148.17", "se")
	fmt.Printf("Q: query %v\n", query)

	out = Dig(query)

	fmt.Printf("Ignoring domain %s\n", domain)

	return out

}
