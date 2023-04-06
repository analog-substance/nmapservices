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
	portsCache          []PortInfo
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

func getAllPorts() []PortInfo {
	if len(portsCache) == 0 {
		file := getFile()
		lines := grep.FileLineByLine(file, portLineRe)

		for line := range lines {
			matches := portLineRe.FindStringSubmatch(line)

			portNumber, _ := strconv.Atoi(matches[2])
			weight, _ := strconv.ParseFloat(matches[4], 32)
			protocol := matches[3]
			service := matches[1]

			portsCache = append(portsCache, PortInfo{
				Service:  service,
				Protocol: protocol,
				Port:     portNumber,
				Weight:   weight,
			})
		}

		sort.SliceStable(portsCache, func(i, j int) bool {
			return portsCache[i].Weight > portsCache[j].Weight
		})
	}

	return portsCache
}

func Ports(protocol Protocol, services ...string) []PortInfo {
	allPorts := getAllPorts()

	var ports []PortInfo
	for _, port := range allPorts {
		if (strings.EqualFold(port.Protocol, string(protocol)) || protocol == Any) &&
			port.ServiceMatch(services...) {
			ports = append(ports, port)
		}
	}

	return ports
}

func TCPPorts(services ...string) []PortInfo {
	return Ports(TCP, services...)
}

func TopTCPPorts(num int, services ...string) []PortInfo {
	ports := Ports(TCP, services...)
	if num > len(ports) {
		return ports
	}

	return ports[:num]
}

func UDPPorts(services ...string) []PortInfo {
	return Ports(UDP, services...)
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
