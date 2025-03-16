package dig

import (
	//"fmt"
	"encoding/json"
	//"fmt"
	"log"
	//"os/exec"
	"os"
)

func OSdig(domain string) string {
	//path := "/Users/tolvmannen/dig.sh "
	var out string
	var jblob []DigJson
	//jsonstring, err := exec.Command("/bin/bash", path, domain).Output()
	jsonstring, err := os.ReadFile("./mockup/www.iis-mock.json")
	//jsonstring, err := exec.Command(path, domain).Output()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte(jsonstring), &jblob)

	out = ToMermaid(jblob)

	return out
	//return string(out)
}

type RR struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Type  string `json:"type"`
	Ttl   string `json:"ttl"`
	Data  string `json:"data"`
}

type DigJson struct {
	Id             uint32   `json:"id"`
	Opcode         string   `json:"opcode"`
	Status         string   `json:"status"`
	Flags          []string `json:"flags"`
	Query_num      uint32   `json:"query_num"`
	Answer_num     uint32   `json:"answer_num"`
	Authority_num  uint32   `json:"authority_num"`
	Additional_num uint32   `json:"additional_num"`
	// Opt_pseudusection
	Question   RR     `json:"question"`
	Answer     []RR   `json:"answer"`
	Authority  []RR   `json:"authority"`
	Additional []RR   `json:"additional"`
	Query_time uint32 `json:"query_time"`
	Server     string `json:"server"`
	When       string `json:"when"`
	Rcvd       uint32 `json:"rcvd"`
	When_epoch uint32 `json:"when_epoch"`
}

func ToMermaid(dj []DigJson) string {

	var graph Graph
	var out string

	var lastnode string
	for _, lv := range dj {
		if lastnode != "" {
			line := Line{
				From: lastnode,
				To:   lv.Question.Name,
			}
			graph.Lines = append(graph.Lines, line)

		}
		node := Node{
			Name: lv.Question.Name,
			NS:   lv.Server,
		}
		graph.Nodes = append(graph.Nodes, node)

		lastnode = lv.Question.Name

	}

	out = graph.ToCode()

	return out

}

type Node struct {
	Name string
	NS   string
}

type Line struct {
	From string
	To   string
}

type Graph struct {
	Nodes []Node
	Lines []Line
}

func (g *Graph) ToCode() string {
	var out string
	for _, node := range g.Nodes {
		out += node.Name + "@{shape: rect, label \"" + node.NS + "\"}"
	}
	for _, line := range g.Lines {
		out += line.From + " --> " + line.To
	}

	return out
}
