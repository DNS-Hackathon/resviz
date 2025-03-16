package dig

import (
	"encoding/json"
	"regexp"
	"strings"

	"fmt"
	"os/exec"
)

type DNSRecord struct {
	ID       int      `json:"id"`
	Opcode   string   `json:"opcode"`
	Status   string   `json:"status"`
	Flags    []string `json:"flags"`
	Question struct {
		Name  string `json:"name"`
		Class string `json:"class"`
		Type  string `json:"type"`
	} `json:"question"`
	Answer     []DNSAnswer `json:"answer"`
	Authority  []DNSAnswer `json:"authority"`
	Additional []DNSAnswer `json:"additional"`
	Server     string      `json:"server"`
}

type DNSAnswer struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Type  string `json:"type"`
	TTL   int    `json:"ttl"`
	Data  string `json:"data"`
}

type DNSAddress struct {
	IP       string
	Port     string
	Hostname string
	Protocol string
}

func ParseServer(input string) (*DNSAddress, error) {
	re := regexp.MustCompile(`(?P<ip>[0-9a-fA-F:.]+)#(?P<port>\d+)\((?P<hostname>[^)]+)\) \((?P<protocol>[^)]+)\)`)
	match := re.FindStringSubmatch(input)

	if match == nil {
		return nil, fmt.Errorf("invalid input format")
	}

	result := &DNSAddress{
		IP:       match[1],
		Port:     match[2],
		Hostname: match[3],
		Protocol: match[4],
	}

	return result, nil
}

var hitServers []string

func OSdig(domain string) string {
	script := "./dig/dig.sh"
	jsonData, err := exec.Command("/bin/sh", script, domain).Output()
	if err != nil {
		return fmt.Sprintf("Error running dig script:%s\n", err)
	}

	var dnsRecords []DNSRecord
	err = json.Unmarshal([]byte(jsonData), &dnsRecords)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON:%s\n", err)
	}

	var mermaidBuilder strings.Builder
	mermaidBuilder.WriteString("graph TD;\n")

	for _, record := range dnsRecords {
		server := fmt.Sprintf("%s\n", record.Server)

		parsed, err := ParseServer(server)
		if err != nil {
			return fmt.Sprintf("Error:", err)
		}

		hitServers = append(hitServers, parsed.Hostname)
		//mermaidBuilder.WriteString(parsed.Hostname)
		//mermaidBuilder.WriteString("\n")
	}

	hitServers = hitServers[1:]
	for _, server := range hitServers {
		mermaidBuilder.WriteString(server + "@{ shape: rect, label: \"" + server + "\"}")
		mermaidBuilder.WriteString("\n")
	}

	mermaidBuilder.WriteString("\n\n")

	for i := 0; i < len(hitServers)-1; i++ {
		mermaidBuilder.WriteString(fmt.Sprintf("%s --> %s\n", hitServers[i], hitServers[i+1]))
	}

	// return the mermaid string object
	return fmt.Sprintf(mermaidBuilder.String())
	// return out
}
