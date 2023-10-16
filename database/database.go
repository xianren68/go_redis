package database

import (
	"fmt"
	"go_redis/config"
	"go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/lib/logger"
	"go_redis/resp/reply"
	"runtime/debug"
	"strconv"
	"strings"
)

/* 测试数据
*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
*2\r\n$3\r\nget\r\n$4\r\nkey1\r\n
*3\r\n$6\r\nrename\r\n$3\r\nkey\r\n$4\r\nkey1\r\n
*2\r\n$3\r\nget\r\n$4\r\nkey1\r\n
 */

// DataBase 数据库
type DataBase struct {
	dbSet []*DB
}

// Exec 执行数据库命令
func (d *DataBase) Exec(connection resp.Connection, line database.CmdLine) resp.Reply {
	// 错误恢复
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
		}
	}()
	// 判断是什么类型的命令
	cName := strings.ToLower(string(line[0]))
	// 选择db的命令
	if cName == "select" {
		if len(line) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(connection, line[1:], d)
	}
	db := d.dbSet[connection.GetDBIndex()]
	return db.Exec(connection, line)
}

func (d *DataBase) Close() {

}

func (d *DataBase) AfterClientClose(connection resp.Connection) {
}

func NewDataBase() *DataBase {
	// 默认db数量为16个
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	db := &DataBase{
		dbSet: make([]*DB, config.Properties.Databases),
	}
	for i := range db.dbSet {
		db.dbSet[i] = MakeDB()
	}
	return db
}
func execSelect(c resp.Connection, cmd [][]byte, d *DataBase) resp.Reply {
	dbIndex, err := strconv.Atoi(string(cmd[0]))
	if err != nil {
		return reply.MakeStandErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(d.dbSet) {
		return reply.MakeStandErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
