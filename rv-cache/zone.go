package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
	cmap "github.com/orcaman/concurrent-map/v2"
	//"github.com/orcaman/concurrent-map/v2"
	//"github.com/onsi/ginkgo/v2/ginkgo/labels"
)

// github.com/google/go-cmp

// Liat of all found nameservers and their IP-addresses
// Limit use to parents of (functionally) second level domains
// The name server Name is the Key.

//type ZONElist map[string]Zone

type NSset map[string]NS

type Zone struct {
	NS         NSset      `json:"NS"`
	ParentData ParentData `json:"ParentData"`
	Primed     bool       `json:"Primed"`
}

type NS struct {
	IP4 []string `json:"IP4"`
	IP6 []string `json:"IP6"`
	SOA uint32   `json:"SOA"`
}

type ParentData struct {
	Delegation NSset  `json:"NS"`
	DS         string `json:"DS"`
}

func (z Zone) Print() {
	jz, err := json.Marshal(z)
	if err != nil {
		fmt.Printf("Warning: unable to unmarshal(%s)\n", err.Error())
	}
	fmt.Printf("\n%+v\n", string(jz))
}

// preload zone data from json file
// mainly for root-hints and se-hints
func (z *Zone) Preload(file string) {

	js, err := os.ReadFile("hints/" + file)
	if err != nil {
		fmt.Printf("Warning: unable to load hint file: %s (%s)\n", file, err.Error())
	}
	json.Unmarshal(js, z)
	if err != nil {
		fmt.Printf("Warning: unable to unmarshal(%s)\n", err.Error())
	}
}

// Walk the DNS tree from root to zone
// Follows Method Get-Parent-NS-IP
func GetParentNSIP(zcache *cmap.ConcurrentMap[string, Zone], childzone string) (string, error) {

	parentzone := ""

	// Zone is not allowed be empty or "ROOT"
	if childzone == "" {
		err := errors.New("Error: Zone name is empty string")
		return parentzone, err
	}

	if childzone == "." {
		err := errors.New("Error: Zone is ROOT (.)")
		return parentzone, err
	}

	// strip one label to get parent zone

	labels := dns.SplitDomainName(childzone)
	if len(labels) == 1 {
		parentzone = "."
	} else {
		parentzone = strings.Join(labels[1:], ".")
	}

	if _, ok := zcache.Get(parentzone); ok {
		fmt.Printf("in found parent for %s = %s\n", childzone, parentzone)
	} else {
		fmt.Printf("dig NS for %s @ %s\n", childzone, parentzone)

	}

	return parentzone, nil

}

/*
 */
