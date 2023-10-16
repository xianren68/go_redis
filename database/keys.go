// Package database 数据库keys指令集
package database

import (
	"go_redis/interface/resp"
	"go_redis/lib/wildcard"
	"go_redis/resp/reply"
)

// 从数据库中删除key,可以是多个
func execDel(db *DB, cmd [][]byte) resp.Reply {
	deleted := 0
	for _, bytes := range cmd {
		deleted += db.Remove(string(bytes))
	}
	return reply.MakeIntReply(int64(deleted))
}

// 检查输入的键中有几个存在
func execExists(db *DB, cmd [][]byte) resp.Reply {
	result := 0
	for _, bytes := range cmd {
		if _, exist := db.GetEntity(string(bytes)); exist {
			result++
		}
	}
	return reply.MakeIntReply(int64(result))
}

// 清空数据库中所有数据
func execFlushDB(db *DB, cmd [][]byte) resp.Reply {
	db.Flush()
	return reply.MakeOkReply()
}

// 返回键对应值的类型
func execType(db *DB, cmd [][]byte) resp.Reply {
	// 拿到对应的键
	key := string(cmd[0])
	val, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeStatusReply("none")
	}
	switch val.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}
	return &reply.UnKnowErrReply{}
}

// 将键重命名
//
//	k1 -> v1 => k2 -> v1
func execRename(db *DB, cmd [][]byte) resp.Reply {
	k1 := string(cmd[0])
	k2 := string(cmd[1])
	val, exist := db.GetEntity(k1)
	if !exist {
		return reply.MakeStandErrReply("no such key")
	}
	// 将旧的删除，新的添加
	db.Remove(k1)
	db.PutEntity(k2, val)
	return reply.MakeOkReply()
}

// 将键重命名，只有新键不存在时有效，防止覆盖原来的值
func execRenameNx(db *DB, cmd [][]byte) resp.Reply {
	k1 := string(cmd[0])
	k2 := string(cmd[1])
	_, exist := db.GetEntity(k2)
	if exist {
		return reply.MakeIntReply(0)
	}
	val, exist := db.GetEntity(k1)
	if !exist {
		return reply.MakeStandErrReply("no such key")
	}
	// 将旧的删除，新的添加
	db.Remove(k1)
	db.PutEntity(k2, val)
	return reply.MakeIntReply(1)
}

// 返回符合通配符选择器的键
func execKeys(db *DB, cmd [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(cmd[0]))
	res := make([][]byte, db.dict.Len())
	// 遍历所有键，找出符合的
	db.dict.ForEach(func(a1 any, a2 any) bool {
		key := a1.(string)
		if pattern.IsMatch(key) {
			res = append(res, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(res)
}

// 注册命令
func init() {
	RegisterCmd("Del", execDel, -2)
	RegisterCmd("Exists", execExists, -2)
	RegisterCmd("Keys", execKeys, 2)
	RegisterCmd("FlushDB", execFlushDB, -1)
	RegisterCmd("Type", execType, 2)
	RegisterCmd("Rename", execRename, 3)
	RegisterCmd("RenameNx", execRenameNx, 3)
}
