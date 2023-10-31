package naabu

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"time"
)

func (naa *Naabu) Type() string {
	return typeof
}

func (naa *Naabu) Name() string {
	return naa.cfg.name
}

func (naa *Naabu) Close() error {
	if naa.task != nil {
		naa.task.close()
	}

	if naa.queue != nil {
		naa.queue.Close()
	}
	naa.UndoDefine()

	return nil
}

func (naa *Naabu) Worker() {
	naa.queue.Subscribe(naa)
}

func (naa *Naabu) Start() error {
	for i := 0; i < naa.cfg.thread; i++ {
		go naa.Worker()
	}
	naa.V(time.Now(), lua.VTRun)
	return nil
}

func (naa *Naabu) pipeL(L *lua.LState) int {
	naa.cfg.Chains.CheckMany(L)
	return 0
}

func (naa *Naabu) NewTaskL(L *lua.LState) int {
	if naa.IsWorking() {
		L.RaiseError("%+v\n", "scan task running")
		return 0
	}

	target := L.CheckString(1)
	task := naa.NewTask([]string{target})
	L.Push(task)
	return 1
}

func (naa *Naabu) defineL(L *lua.LState) int {
	naa.Define()
	return 0
}

func (naa *Naabu) startL(L *lua.LState) int {
	xEnv.Start(L, naa).
		Err(func(err error) { L.RaiseError("%v", err) }).
		From(L.CodeVM()).
		Do()
	return 0
}

func (naa *Naabu) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "start":
		return lua.NewFunction(naa.startL)
	case "pipe":
		return lua.NewFunction(naa.pipeL)

	case "task":
		return lua.NewFunction(naa.NewTaskL)

	case "define":
		return lua.NewFunction(naa.defineL)
	default:
		//todo
	}

	return lua.LNil
}
