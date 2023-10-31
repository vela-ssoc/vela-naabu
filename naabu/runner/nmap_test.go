package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vela-ssoc/vela-naabu/naabu/port"
	"github.com/vela-ssoc/vela-naabu/naabu/protocol"
	"github.com/vela-ssoc/vela-naabu/naabu/result"
	"github.com/vela-ssoc/vela-naabu/naabu/scan"
)

func TestHandleNmap(t *testing.T) {
	// just attempt to start nmap
	var r Runner
	r.options = &Options{}
	// nmap with empty cli shouldn't trigger any error
	res := result.NewResult()
	r.scanner = &scan.Scanner{}
	r.scanner.ScanResults = res
	assert.Nil(t, r.handleNmap())
	// nmap syntax error (this test might fail if nmap is not installed on the box)
	assert.Nil(t, r.handleNmap())
	r.scanner.ScanResults.SetPorts("127.0.0.1", []*port.Port{{Port: 8080, Protocol: protocol.TCP}})
	assert.Nil(t, r.handleNmap())
}
