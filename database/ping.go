package database

import (
	"go_redis/interface/resp"
	"go_redis/resp/reply"
)

func execPing(db *DB, cmd [][]byte) resp.Reply {
	if len(cmd) == 0 {
		return reply.MakeOkReply()
	} else if len(cmd) == 1 {
		return reply.MakeStatusReply(string(cmd[0]))
	} else {
		return reply.MakeStandErrReply("ERR wrong number of arguments for 'ping' command")
	}
}
