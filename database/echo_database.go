package database

import (
	databaseface "go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

// 执行命令

func (e *EchoDatabase) Exec(connection resp.Connection, line databaseface.CmdLine) resp.Reply {
	return reply.MakeMultiBulkReply(line)
}

func (e *EchoDatabase) Close() {

}

func (e *EchoDatabase) AfterClientClose(connection resp.Connection) {

}
