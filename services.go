package nmapservices

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/analog-substance/fileutil/grep"
	"github.com/analog-substance/nmapservices/internal/static"
)

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
	Any Protocol = ""
)

var (
	defaultServicesPath string         = filepath.FromSlash("/usr/share/nmap/nmap-services")
	portLineRe          *regexp.Regexp = regexp.MustCompile(`^([^\t]+)\t([0-9]+)/(tcp|udp)\t([0-9\.]+)`)
	//portsCache          []PortInfo // Will probably want to have a cache of all the ports for efficiency at some point
)

type PortInfo struct {
	Service  string
	Protocol string
	Port     int
	Weight   float64
}

func (p PortInfo) ServiceMatch(services ...string) bool {
	if len(services) == 0 {
		return true
	}

	for _, service := range services {
		if strings.Contains(p.Service, service) {
			return true
		}
	}
	return false
}

func getFile() fs.File {
	var file fs.File

	file, err := os.Open(defaultServicesPath)
	if err != nil {
		file, _ = static.Files.Open("nmap-services")
	}
	return file
}

func Ports(protocol Protocol, services ...string) []PortInfo {
	file := getFile()
	lines := grep.FileLineByLine(file, portLineRe)

	var ports []PortInfo
	for line := range lines {
		matches := portLineRe.FindStringSubmatch(line)

		portNumber, _ := strconv.Atoi(matches[2])
		weight, _ := strconv.ParseFloat(matches[4], 32)
		proto := matches[3]
		service := matches[1]

		portInfo := PortInfo{
			Service:  service,
			Protocol: proto,
			Port:     portNumber,
			Weight:   weight,
		}

		if (strings.EqualFold(proto, string(protocol)) || protocol == Any) &&
			portInfo.ServiceMatch(services...) {
			ports = append(ports, portInfo)
		}
	}

	sort.SliceStable(ports, func(i, j int) bool {
		return ports[i].Weight > ports[j].Weight
	})

	return ports
}

func TopTCPPorts(num int, services ...string) []PortInfo {
	ports := Ports(TCP, services...)
	if num > len(ports) {
		return ports
	}

	return ports[:num]
}

func TopUDPPorts(num int, services ...string) []PortInfo {
	ports := Ports(UDP, services...)
	if num > len(ports) {
		return ports
	}

	return ports[:num]
}

func init() {
	if runtime.GOOS == "windows" {
		defaultServicesPath = `` // Currently don't know the default location
	}
}
