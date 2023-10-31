package naabu

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

func NewNaabuL(L *lua.LState) int {
	cfg := NewConfig(L)
	vda := L.NewVelaData(cfg.name, typeof) //判断出 当前code 是否有相同的对象 名字和类型
	if vda.IsNil() {
		vda.Set(NewNaabu(cfg))
		L.Push(vda)
	} else {
		old := vda.Data.(*Naabu)
		old.cfg = cfg
		L.Push(vda)
	}
	return 1

}

/*
	local s = vela.scanner{
	    name   = "scanner",
		thread = 10,
		//指纹库
		task   = 1,
		report = vela.report()
	}

	s.pipe(function(host)
		//host.peer host.port host.banner  host.app  host.http_raw
	end)

	s.start()

	s.new("172.31.61.0/24").ov().ss().start() -- { taskid:xxxx }
*/

func WithEnv(env vela.Environment) {
	xEnv = env
	tab := lua.NewUserKV()
	xEnv.Set("naabu", lua.NewExport("vela.naabu.export", lua.WithTable(tab), lua.WithFunc(NewNaabuL)))
}
