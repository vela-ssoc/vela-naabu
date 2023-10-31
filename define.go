package naabu

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-naabu/naabu/runner"
)

type TaskParam struct {
	ID                 string   `json:"id"`
	Host               []string `json:"host"`
	Mode               string   `json:"mode"`
	Ports              string   `json:"ports"`
	AllPort            bool     `json:"all_port"`
	Rate               int      `json:"rate"`
	Retries            int      `json:"retries"`
	Ping               bool     `json:"ping"`
	Verbose            bool     `json:"verbose"`
	ServiceDiscovery   bool     `json:"service_discovery"`
	SkipHostDiscovery  bool     `json:"skip_host_discovery"`
	DisableUpdateCheck bool     `json:"disable_update_check"`
	ServiceVersion     bool     `json:"service_version"`
}

func NewTaskParam() *TaskParam {
	return &TaskParam{
		Verbose: false,
		Mode:    "sn",
		// ScanAllIPS: true,
		// Ports: "80-10000",
		Retries:            2,
		Rate:               256, //wg thread worker goroutine
		Ping:               false,
		ServiceDiscovery:   true,
		SkipHostDiscovery:  true,
		DisableUpdateCheck: true,
		ServiceVersion:     true,
	}

}

func (tp *TaskParam) Clone(option *runner.Options) {
	option.Host = tp.Host
	option.Verbose = tp.Verbose
	option.ScanType = tp.Mode
	option.ScanAllIPS = tp.AllPort
	option.Ports = tp.Ports
	option.Retries = tp.Retries
	option.Rate = tp.Rate //wg thread worker goroutine
	option.Ping = tp.Ping
	option.ServiceDiscovery = tp.ServiceDiscovery
	option.SkipHostDiscovery = tp.SkipHostDiscovery
	option.DisableUpdateCheck = tp.DisableUpdateCheck
	option.ServiceVersion = tp.ServiceVersion
}

func (naa *Naabu) TaskPath() string {
	return fmt.Sprintf("/api/v1/arr/agent/lua/naabu/%s/task", naa.Name())
}

func (naa *Naabu) StatusPath() string {
	return fmt.Sprintf("/api/v1/arr/agent/lua/naabu/%s/status", naa.Name())
}

func (naa *Naabu) TaskHandle(ctx *fasthttp.RequestCtx) error {
	param := NewTaskParam()
	body := ctx.PostBody()
	if len(body) == 0 {
		return fmt.Errorf("task handle fail got empty")
	}

	err := json.Unmarshal(body, &param)
	if err != nil {
		return err
	}

	task := naa.NewTask(param.Host)
	param.Clone(task.Option)
	ctx.Write(task.info())
	go task.GenRun()
	return nil
}

func (naa *Naabu) StatusHandle(ctx *fasthttp.RequestCtx) error {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("task_id", naa.TaskID())
	enc.KV("task_status", naa.TaskStatus())
	enc.End("}")
	ctx.Write(enc.Bytes())
	return nil
}

func (naa *Naabu) Define() {
	r := xEnv.R()
	r.POST(naa.TaskPath(), xEnv.Then(naa.TaskHandle))
	r.GET(naa.StatusPath(), xEnv.Then(naa.StatusHandle))
}

func (naa *Naabu) UndoDefine() {
	r := xEnv.R()
	r.Undo(fasthttp.MethodPost, naa.TaskPath())
	r.Undo(fasthttp.MethodGet, naa.StatusPath())
}
