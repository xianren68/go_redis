package database

import (
	"go_redis/interface/resp"
)

type CmdLine = [][]byte
type Database interface {
	Exec(resp.Connection, CmdLine) resp.Reply
	Close()
	AfterClientClose(resp.Connection)
}
type DataEntity struct {
	Data any
}
