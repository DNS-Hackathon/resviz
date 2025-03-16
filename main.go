package main

import (
	//"encoding"
	"fmt"
	//"os"
	cache "resviz/rv-cache"
	engine "resviz/rv-engine"
	//dig "resviz/dig"
	//"encoding/json"
)

// set up Global Cache
var GC cache.Cache

func main() {

	debug := false

	GC.Preload()

	if debug {
		fmt.Printf("%+v\n", GC.Zones)
	}

	engine.Run()

}
