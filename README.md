# vela-Naabu
安全资产探测插件

```lua
local naa = vela.naabu{
  name = "scanner",
  finger = {timeout = 500 , udp = false , fast = false},
}
naa.start()

local es = vela.elastic.default("vela-naabu-%s" , "$day")
naa.pipe(function(host)
  es.send(host)
end)

-- 外部API
naa.define()
```
