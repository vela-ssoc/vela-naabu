package naabu

import (
	"context"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-naabu/naabu/port"
	"github.com/vela-ssoc/vela-naabu/naabu/result"
	"github.com/vela-ssoc/vela-naabu/naabu/runner"
	"time"
)

type Dispatch interface {
	End(*Task)
	OnResult(*result.HostResult)
	Callback(string, *port.Port)
	Catch(error)
}

type Task struct {
	ID       string
	ctime    time.Time
	ctx      context.Context
	cancel   context.CancelFunc
	Option   *runner.Options
	Runner   *runner.Runner
	dispatch Dispatch
}

func (t *Task) String() string                         { return "" }
func (t *Task) Type() lua.LValueType                   { return lua.LTObject }
func (t *Task) AssertFloat64() (float64, bool)         { return 0, false }
func (t *Task) AssertString() (string, bool)           { return "", false }
func (t *Task) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (t *Task) Peek() lua.LValue                       { return t }

func (t *Task) close() error {
	if t.cancel == nil {
		return nil
	}
	t.cancel()
	return nil
}

func (t *Task) info() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("id", t.ID)
	enc.KV("status", "working")
	enc.KV("ctime", t.ctime)
	enc.End("}")
	return enc.Bytes()
}

func (t *Task) GenRun() {
	if t.dispatch == nil {
		xEnv.Errorf("%s dispatch got nil")
		return
	}

	t.Option.OnResult = t.dispatch.OnResult
	t.Option.Callback = t.dispatch.Callback

	r, err := runner.NewRunner(t.Option)
	if err != nil {
		t.dispatch.Catch(err)
		return
	}
	t.Runner = r

	defer t.close()
	r.RunEnumeration(t.ctx)
	t.dispatch.End(t)
}
