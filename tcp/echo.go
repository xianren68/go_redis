package tcp

import (
	"bufio"
	"context"
	"go_redis/lib/logger"
	"go_redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// 客户端连接
type EchoClient struct {
	Conn net.Conn
	// 超时时间(等待组)
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	// 设置超时等待(防止业务没做完)
	e.Waiting.WaitWithTimeout(10 * time.Second)
	// 关闭连接
	_ = e.Conn.Close()
	return nil

}

// 业务引擎
type EchoHandler struct {
	// 记录连接的客户端
	activeConn sync.Map
	// 业务引擎是否关闭
	closing atomic.Bool
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if e.closing.Load() {
		// 业务引擎关闭，新的连接直接断开
		_ = conn.Close()
	}
	// 包装连接
	client := &EchoClient{
		Conn: conn,
	}
	// 记录客户端连接
	e.activeConn.Store(client, struct{}{})
	// 执行业务
	reader := bufio.NewReader(conn)
	for {
		// 监听客户端发送的信息
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// 客户端关闭
				logger.Info("connect close")
				// 将客户端移出连接列表
				e.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return

		}
		// 等待组，防止在做业务时被关闭
		client.Waiting.Add(1)
		// 将数据返回给客户端
		conn.Write([]byte(msg))
		client.Waiting.Done()

	}

}
func (e *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	// 修改状态
	e.closing.Store(false)
	// 将所有的连接断开
	e.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		client.Close()
		return true
	})
	return nil
}
