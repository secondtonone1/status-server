package redis

import (
	"io/ioutil"
	"status-server/logging"
)

var (
	StatusLua LuaScript
)

type LuaScript struct {
	LuaPath string
	Script  string
	Sha     string
}

/**
redis加载lua脚本
luaPath：文件存放路径
isStoreScript: 是否缓存到redis是。解决每次调用传脚本
*/
func FromResource(luaPath string, isStoreScript bool) LuaScript {
	//读取脚本文件
	scriptByte, err2 := ioutil.ReadFile(luaPath)
	if err2 != nil {
		panic("read file error [" + luaPath + "]" + err2.Error())
	}

	//加载到redis中进行缓存
	script := string(scriptByte)
	luaScript := LuaScript{LuaPath: luaPath, Script: script}

	//判断是否需要保存到redis上
	if isStoreScript {
		err := luaScript.storeScript(&luaScript)
		if err != nil {
			panic(err.Error())
		}
	}
	return luaScript
}

//redis保存lua脚本
func (obj LuaScript) storeScript(l *LuaScript) error {
	r, e := ScriptLoad(obj.Script)
	if e != nil {
		logging.Logger.Errorf("redis_util.FromResource加载redis lua脚本异常！path[%s],script[%s] ,err[%s]", obj.LuaPath, obj.Script, e.Error())
		return e
	} else {
		l.Sha = r
		return nil
	}
}

//执行lua脚本
func (obj LuaScript) EvalSha(keys []string, args ...interface{}) (interface{}, error) {
	r, e := EvalSha(obj.Sha, keys, args...)
	if e != nil {
		e2 := obj.storeScript(&obj)
		if e2 == nil {
			return EvalSha(obj.Sha, keys, args...)
		}
	}
	return r, e
}

//执行lua脚本
func (obj LuaScript) Eval(keys []string, args ...interface{}) (interface{}, error) {
	r, e := Eval(obj.Script, keys, args...)
	if e != nil {
		return  Eval(obj.Script, keys, args...)
	}
	return r, e
}
