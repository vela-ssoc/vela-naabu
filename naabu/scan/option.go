package scan

import (
	"github.com/vela-ssoc/vela-naabu/naabu/port"
	"time"
)

// Options of the scan
type Options struct {
	Callback      func(string, *port.Port)
	Timeout       time.Duration
	Retries       int
	Rate          int
	PortThreshold int
	Debug         bool
	ExcludeCdn    bool
	OutputCdn     bool
	ExcludedIps   []string
	Proxy         string
	ProxyAuth     string
	Stream        bool
}
