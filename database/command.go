package database

import "strings"

// 所有的指令列表
var cmdTable = map[string]*command{}

// 数据库指令
type command struct {
	// 具体执行
	executor ExecFunc
	// 这个指令需要的参数
	arity int
}

// RegisterCmd 注册数据库指令
func RegisterCmd(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmd := &command{executor, arity}
	cmdTable[name] = cmd
}
