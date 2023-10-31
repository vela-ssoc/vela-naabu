package naabu

import "github.com/vela-ssoc/vela-kit/lua"

func (t *Task) nmapL(L *lua.LState) int {
	t.Option.NmapCLI = L.CheckString(1)
	L.Push(t)
	return 1
}

func (t *Task) ConnL(L *lua.LState) int {
	t.Option.ScanType = "c"
	L.Push(t)
	return 1
}

func (t *Task) SynL(L *lua.LState) int {
	t.Option.ScanType = "s"
	L.Push(t)
	return 1
}

func (t *Task) AllL(L *lua.LState) int {
	t.Option.ScanAllIPS = L.IsTrue(1)
	L.Push(t)
	return 1
}

func (t *Task) runL(L *lua.LState) int {
	go t.GenRun()
	return 0
}

func (t *Task) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "nmap":
		return lua.NewFunction(t.nmapL)
	case "syn":
		return lua.NewFunction(t.SynL)
	case "conn":
		return lua.NewFunction(t.ConnL)
	case "all":
		return lua.NewFunction(t.AllL)
	case "run":
		return lua.NewFunction(t.runL)
	default:
		return lua.LNil
	}

}
