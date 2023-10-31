package naabu

import (
	"context"
	"github.com/google/uuid"
	"github.com/vela-ssoc/vela-kit/chanutil"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-naabu/naabu/port"
	"github.com/vela-ssoc/vela-naabu/naabu/result"
	"github.com/vela-ssoc/vela-naabu/naabu/runner"
	"reflect"
	"sync/atomic"
	"time"
)

var typeof = reflect.TypeOf((*Naabu)(nil)).String()

const (
	Idle uint32 = iota + 1
	Working
)

type Naabu struct {
	lua.SuperVelaData
	Status uint32
	cfg    *Config
	task   *Task
	queue  *chanutil.Queue
}

func (naa *Naabu) IsWorking() bool {
	return atomic.LoadUint32(&naa.Status) == Working
}

func (naa *Naabu) TaskID() string {
	if naa.task == nil {
		return ""
	}

	return naa.task.ID
}

func (naa *Naabu) TaskStatus() string {
	if naa.task == nil {
		return "idle"
	}
	return "working"
}

func (naa *Naabu) Exception(err error) {
	xEnv.Errorf("naabu handle fail %v", err)
}

func (naa *Naabu) Catch(err error) {
	xEnv.Errorf("naabu handle fail %v", err)
}

func (naa *Naabu) OnResult(r *result.HostResult) {

}

func (naa *Naabu) handle(h *Host) {
	naa.cfg.Chains.Do(h, naa.cfg.co, func(err error) {
		naa.Exception(err)
	})
}

func (naa *Naabu) End(t *Task) {
	atomic.StoreUint32(&naa.Status, Idle)
	naa.task = nil
}

func (naa *Naabu) Callback(ip string, p *port.Port) {
	if e := naa.queue.Publish(Tx{IP: ip, Port: p}); e != nil {
		xEnv.Errorf("naabu callback fail %v", e)
	}
}

func (naa *Naabu) NewTask(host []string) *Task {
	option := &runner.Options{
		Host:     host,
		Verbose:  false,
		ScanType: "sn",
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

	ctx, cancel := context.WithCancel(xEnv.Context())
	naa.task = &Task{
		ctx:      ctx,
		cancel:   cancel,
		ID:       uuid.New().String(),
		Option:   option,
		dispatch: naa,
		ctime:    time.Now(),
	}

	return naa.task
}

func NewNaabu(cfg *Config) *Naabu {
	naa := &Naabu{
		cfg:    cfg,
		Status: Idle,
		queue:  chanutil.NewQueue(0),
	}

	return naa
}
