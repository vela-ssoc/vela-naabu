package naabu

import (
	"github.com/vela-ssoc/vela-naabu/fingerprintx/plugins"
	"github.com/vela-ssoc/vela-naabu/fingerprintx/scan"
	"github.com/vela-ssoc/vela-naabu/naabu/protocol"
	"net/netip"
)

// OnMessage naabu queue message
func (naa *Naabu) OnMessage(v interface{}) {
	tx, ok := v.(Tx)
	if !ok {
		return
	}

	fn := func() { //default
		naa.handle(&Host{
			IP:        tx.IP,
			Port:      tx.Port.Port,
			Transport: tx.Port.Protocol.String(),
			TLS:       tx.Port.TLS,
		})
	}
	if naa.cfg.FingerDisable {
		fn()
		return
	}

	addr, err := netip.ParseAddr(tx.IP)
	if err != nil {
		fn()
		return
	}

	target := plugins.Target{
		Address: netip.AddrPortFrom(addr, uint16(tx.Port.Port)),
		Host:    "localhost",
	}

	fx := naa.cfg.Finger()

	switch tx.Port.Protocol {
	case protocol.UDP:
		fx.UDP = true
	case protocol.TCP:
		fx.UDP = false
	case protocol.ARP:
		fx.UDP = true
	}

	srv, err := scan.Do(target, fx)
	if err != nil {
		xEnv.Errorf("callback fingerprint fail %v", err)
		fn()
		return
	}

	if srv == nil {
		fn()
		return
	}

	h := &Host{
		IP:        srv.IP,
		Port:      srv.Port,
		Protocol:  srv.Protocol,
		TLS:       srv.TLS,
		Transport: srv.Transport,
		Version:   srv.Version,
		Banner:    srv.Raw,
	}
	naa.handle(h)
}

func (naa *Naabu) OnClose() {

}
