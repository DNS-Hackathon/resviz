package cache

import (
	"errors"
	//"fmt"
	"github.com/miekg/dns"
	"github.com/orcaman/concurrent-map/v2"
	"strings"
)

//*cmap.ConcurrentMap[string, <Type>]

type Cache struct {
	Zones   cmap.ConcurrentMap[string, Zone]
	NScache cmap.ConcurrentMap[string, NSset] // For caching Authorativve servers in bulktest
}

func (c *Cache) Preload() {
	c.Zones = cmap.New[Zone]()

	var root Zone
	root.Preload("root-hints.json")
	c.Zones.Set(".", root)

	var se Zone
	se.Preload("se-hints.json")
	c.Zones.Set("se", se)

}

func (c *Cache) GetZone(z string) (Zone, error) {

	if z == "" {
		err := errors.New("Error: Zone name is empty string")
		return Zone{}, err
	}

	// Check cache for zone
	if zone, ok := c.Zones.Get(z); ok {
		return zone, nil
	}

	// If not in cache, create new zone
	zone := Zone{}

	// Find NS of parent domain to query for delegation.
	labels := dns.SplitDomainName(z)
	var parent string
	if len(labels) == 1 {
		parent = "."
	} else {
		parent = strings.Join(labels[1:], ".")
	}

	// Get parent Zone from cache (recurive, fuck yeah)
	pz, err := c.GetZone(parent)

	if err != nil {
		return zone, err
	}

	// query for delegation info
	GetDeleg(z, pz.NS)

	return zone, nil

}

func GetDeleg(zone string, parent NSset) (ParentData, error) {

	pd := ParentData{}

	return pd, nil

}
