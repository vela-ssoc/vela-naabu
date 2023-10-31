package naabu

import (
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/strutil"
)

type Host struct {
	IP        string          `json:"ip"`
	Port      int             `json:"port"`
	TLS       bool            `json:"tls"`
	Banner    json.RawMessage `json:"banner"`
	Protocol  string          `json:"protocol"`
	Transport string          `json:"transport"`
	Version   string          `json:"version"`
}

func (h *Host) String() string                         { return strutil.B2S(h.Bytes()) }
func (h *Host) Type() lua.LValueType                   { return lua.LTObject }
func (h *Host) AssertFloat64() (float64, bool)         { return 0, false }
func (h *Host) AssertString() (string, bool)           { return "", false }
func (h *Host) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (h *Host) Peek() lua.LValue                       { return h }

func (h *Host) Bytes() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("ip", h.IP)
	enc.KV("port", h.Port)
	enc.KV("tls", h.TLS)
	enc.KV("protocol", h.Protocol)
	enc.KV("transport", h.Transport)
	enc.KV("version", h.Version)
	enc.Raw("banner", h.Banner)
	enc.End("}")
	return enc.Bytes()
}
