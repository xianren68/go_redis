package resp

import "net"

type Connection interface {
	Write([]byte) error
	Close() error
	// 返回客户端地址

	RemoteAddr() net.Addr
	// 返回db索引
	GetDBIndex() int
	// 选择db
	SelectDB(int)
}
