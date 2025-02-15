package handler

import (
	"goredis/resp"
	"sync"
)

var Handler = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

var SETstore = map[string]string{}
var SETstoreMu = sync.RWMutex{}

var HSETstore = map[string]map[string]string{}
var HSETstoreMu = sync.RWMutex{}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "simplestring", Str: "PONG"}
	}
	return resp.Value{Typ: "simplestring", Str: args[0].Bulk}
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "simpleerror", Str: "SET requires two arguments"}
	}
	key := args[0].Bulk
	val := args[1].Bulk
	SETstoreMu.Lock()
	SETstore[key] = val
	SETstoreMu.Unlock()

	return resp.Value{Typ: "simplestring", Str: "OK"}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "simpleerror", Str: "GET requires one argument"}
	}
	key := args[0].Bulk
	SETstoreMu.RLock()
	val := SETstore[key]
	SETstoreMu.RUnlock()

	return resp.Value{Typ: "bulk", Bulk: val}
}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "simpleerror", Str: "HSET requires three arguments"}
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	val := args[2].Bulk

	HSETstoreMu.Lock()
	if _, ok := HSETstore[hash]; !ok {
		HSETstore[hash] = map[string]string{}
	}
	HSETstore[hash][key] = val
	HSETstoreMu.Unlock()

	return resp.Value{Typ: "simplestring", Str: "OK"}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "simpleerror", Str: "HGET requires two arguments"}
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	SETstoreMu.RLock()
	val, ok := HSETstore[hash][key]
	SETstoreMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: val}
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "simpleerror", Str: "HGETALL needs one argument"}
	}
	hash := args[0].Bulk
	SETstoreMu.RLock()
	valstore, ok := HSETstore[hash]
	SETstoreMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}
	bulks := []resp.Value{}
	for key, value := range valstore {
		bulks = append(bulks, resp.Value{Typ: "bulk", Bulk: key})
		bulks = append(bulks, resp.Value{Typ: "bulk", Bulk: value})
	}
	return resp.Value{Typ: "array", Array: bulks}
}
