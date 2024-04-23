package lua

import (
	"context"
	"errors"
	"log"
	"os"
	"rate-limiter/api/types"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var LuaScriptPath string = "pkg/redis/lua/limiter.lua"
var Sha string

func LoadScript() (*redis.Script, error) {
	luaScript, err := os.ReadFile(LuaScriptPath)
	if err != nil {
		log.Printf("Error reading Lua script: %v", err)
		return nil, err
	}
	return redis.NewScript(string(luaScript)), nil
}

func LoadShaScript(rdb *redis.Client) {
	script, err := os.ReadFile(LuaScriptPath)
	if err != nil {
		log.Printf("Error reading Lua script: %v", err)
		return
	}
	Sha, err = rdb.ScriptLoad(context.Background(), string(script)).Result()
	if err != nil {
		log.Printf("Error loading Lua script SHA: %v", err)
		return
	}

	log.Printf("Lua script loaded successfully with SHA: %v", Sha)
}

func GetSha() string {
	return Sha
}

func ParseResponse(response interface{}) (types.Response, error) {
	responseSlice, ok := response.([]interface{})
	if !ok {
		return types.Response{}, errors.New("error parsing Lua response: expected slice of interface{}")
	}

	successInt, ok := responseSlice[0].(int64)
	if !ok {
		return types.Response{}, errors.New("error parsing Lua response: expected bool for successInt")
	}

	count, ok := responseSlice[1].(int64)
	if !ok {
		return types.Response{}, errors.New("error parsing Lua response: expected int64 for count")
	}

	timeBeforeNextReq, ok := responseSlice[2].(int64)
	if !ok {
		return types.Response{}, errors.New("error parsing Lua response: expected int64 for timeBeforeNextReq")
	}

	success := successInt == 1
	timeBeforeNextReqInSec := timeBeforeNextReq / 1000

	return types.Response{Success: success, Count: count, TimeBeforeNextReq: strconv.FormatInt(timeBeforeNextReqInSec, 10)}, nil
}
