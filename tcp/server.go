package tcp

import (
	"context"
	"go_redis/interface/tcp"
	"go_redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	// 端口
	Address string
}

func ListenAndServerWithSignal(cfg *Config, handler tcp.Handler) error {
	// 创建tcp连接
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	closeChan := make(chan struct{})
	singChan := make(chan os.Signal)
	// 系统信号收集到singchan
	signal.Notify(singChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	// 监听系统信号
	go func() {
		sig := <-singChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			// 发送关闭信号
			closeChan <- struct{}{}
		}
	}()
	ListenAndServer(listener, handler, closeChan)
	return nil
}
func ListenAndServer(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	// 函数退出前关闭端口监听及handler
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	// 监听系统信号
	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	// 防止函数退出时，客户端业务还未完成
	waitGroup := sync.WaitGroup{}
	for {
		// 获取新连接
		conn, err := listener.Accept()
		// 在等待组中+1
		waitGroup.Add(1)
		if err != nil {
			break
		}
		logger.Info("new connect")
		// 每个协程对应一个连接
		go func(conn net.Conn) {
			// 业务完成后-1
			defer waitGroup.Done()
			handler.Handle(ctx, conn)
		}(conn)
	}

}
