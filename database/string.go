package database

import (
	"go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/resp/reply"
)

func (db *DB) getAsString(key string) ([]byte, reply.ErrorReply) {
	value, exist := db.GetEntity(key)
	if !exist {
		return nil, nil
	}
	// 判断值是否能转成字符串
	bytes, ok := value.Data.([]byte)
	if !ok {
		return nil, &reply.WrongTypeErrReply{}
	}
	return bytes, nil
}
func execGet(db *DB, cmd [][]byte) resp.Reply {
	value, err := db.getAsString(string(cmd[0]))
	if err != nil {
		return err
	}
	if value == nil {
		return reply.MakeNullBulkBytes()
	}
	return reply.MakeBulkReply(value)
}
func execSet(db *DB, cmd [][]byte) resp.Reply {
	db.PutEntity(string(cmd[0]), &database.DataEntity{cmd[1]})
	return reply.MakeOkReply()
}

// 添加键值对，当键不存在的时候
func execSetNX(db *DB, cmd [][]byte) resp.Reply {
	return reply.MakeIntReply(int64(db.PutIfAbsent(string(cmd[0]), &database.DataEntity{Data: cmd[1]})))
}

// 获取旧的值，并更新一个新值
func execGetSet(db *DB, cmd [][]byte) resp.Reply {
	key := string(cmd[0])
	old, exist := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: cmd[1]})
	if !exist {
		return reply.MakeNullBulkBytes()
	}
	return reply.MakeBulkReply(old.Data.([]byte))
}

// 获取键对应值的长度
func execStrLen(db *DB, cmd [][]byte) resp.Reply {
	key := string(cmd[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkBytes()
	}
	value := entity.Data.([]byte)
	return reply.MakeIntReply(int64(len(value)))
}
func init() {
	RegisterCmd("Get", execGet, 2)
	RegisterCmd("Set", execSet, -3)
	RegisterCmd("SetNx", execSetNX, 3)
	RegisterCmd("GetSet", execGetSet, 3)
	RegisterCmd("StrLen", execStrLen, 2)
}
