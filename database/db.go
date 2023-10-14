package database

import (
	"go_redis/datastruct/dict"
	"go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/resp/reply"
	"strings"
)

// ExecFunc 数据库执行指令的方法
type ExecFunc = func(*DB, [][]byte) resp.Reply

// DB 数据库
type DB struct {
	// 索引
	index int
	// 数据库接口
	dict dict.Dict
}

func MakeDB() *DB {
	return &DB{
		dict: dict.MakeSyncDict(),
	}
}

// Exec 执行指令
func (db *DB) Exec(connection resp.Connection, cmdline [][]byte) resp.Reply {
	// 转化为小写
	cName := strings.ToLower(string(cmdline[0]))
	// 判断数据库是否存在此指令
	cmd, ok := cmdTable[cName]
	if !ok {
		return reply.MakeStandErrReply("Err unknow command " + cName)
	}
	// 判断命令参数是否符合数量
	if !validateArity(cmd.arity, cmdline) {
		return reply.MakeArgNumErrReply(cName)
	}
	return cmd.executor(db, cmdline[1:])
}

// 判断指令所需参数是否合理
//
// 参数：
//
//	arity: 所需要的参数个数,负数代表不定长参数，以及其最少需要的个数
//	cmdline: 传入的命令
func validateArity(arity int, cmdline [][]byte) bool {
	argNum := len(cmdline)
	if arity < 0 {
		return argNum >= -arity
	}
	return argNum == arity
}

/* ————通用指令———— */

// GetEntity 获取某个键对应的值
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	val, exist := db.dict.Get(key)
	// 不存在
	if !exist {
		return nil, false
	}
	return &database.DataEntity{Data: val}, true
}

// PutEntity 插入键值对
func (db *DB) PutEntity(key string, val *database.DataEntity) int {
	return db.dict.Put(key, val.Data)
}

// PutIfAbsent 插入键值对，当键值不存在的时候
func (db *DB) PutIfAbsent(key string, val *database.DataEntity) int {
	return db.dict.PutIfAbsent(key, val.Data)
}

// PutIfExist 插入键值对，当键值存在的时候(更新)
func (db *DB) PutIfExist(key string, val *database.DataEntity) int {
	return db.dict.PutIfExist(key, val.Data)
}

// Remove 删除
func (db *DB) Remove(key string) int {
	return db.dict.Remove(key)
}

// Removes 删除多个
func (db *DB) Removes(keys ...string) int {
	delete := 0
	for i := 0; i < len(keys); i++ {
		delete += db.Removes(keys[i])
	}
	return delete
}

// Flush 清空数据库
func (db *DB) Flush() {
	db.dict.Clear()
}
