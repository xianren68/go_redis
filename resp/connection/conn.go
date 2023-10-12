package connection

import (
	"go_redis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

// 描述客户端连接的结构体

type Connection struct {
	conn       net.Conn
	waiting    wait.Wait
	mu         sync.Mutex
	selectedDB int
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}
func (c *Connection) Write(bytes []byte) error {
	defer func() {
		c.mu.Unlock()
		c.waiting.Done()
	}()
	if len(bytes) == 0 {
		return nil
	}
	c.mu.Lock()
	c.waiting.Add(1)
	_, err := c.conn.Write(bytes)
	return err
}

func (c *Connection) Close() error {
	// 超时等待
	c.waiting.WaitWithTimeout(10 * time.Second)
	c.conn.Close()
	return nil
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c *Connection) SelectDB(i int) {
	c.selectedDB = i
}
