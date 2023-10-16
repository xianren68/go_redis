package handler

import (
	"context"
	"go_redis/database"
	databaseface "go_redis/interface/database"
	"go_redis/lib/logger"
	"go_redis/resp/connection"
	"go_redis/resp/parse"
	"go_redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

type RespHandler struct {
	// 记录连接的客户端
	activeConn sync.Map
	// 业务引擎是否关闭
	closing atomic.Bool
	db      databaseface.Database
}

func MakeRespHandler() *RespHandler {
	var db databaseface.Database
	db = database.NewDataBase()
	return &RespHandler{
		db: db,
	}
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Load() {
		_ = conn.Close()
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parse.ParseStream(conn)
	for palyLoad := range ch {
		if palyLoad.Err != nil {
			// EOF
			if palyLoad.Err == io.EOF || palyLoad.Err == io.ErrUnexpectedEOF || strings.Contains(palyLoad.Err.Error(),
				"use of closed network connection") {
				r.closeClient(client)
				logger.Info("connection closed" + client.RemoteAddr().String())
				return
			}
			// Protocol err
			res := reply.MakeStandErrReply(palyLoad.Err.Error())
			err := client.Write(res.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed" + client.RemoteAddr().String())
				return
			}
			continue
		} else {
			// 正常逻辑
			val, ok := palyLoad.Data.(*reply.MultiBulkReply)
			if !ok {
				logger.Warn("require multiBulk reply")
				continue
			}
			exec := r.db.Exec(client, val.Args)
			if exec != nil {
				_ = client.Write(exec.ToBytes())
			} else {
				_ = client.Write((&reply.UnKnowErrReply{}).ToBytes())
			}

		}

	}
}

func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Store(true)
	r.activeConn.Range(func(key, value any) bool {
		key.(*connection.Connection).Close()
		return true
	})
	r.db.Close()
	return nil
}

// 关闭某个客户端连接
func (r *RespHandler) closeClient(connection *connection.Connection) {
	connection.Close()
	r.db.AfterClientClose(connection)
	r.activeConn.Delete(connection)
}
