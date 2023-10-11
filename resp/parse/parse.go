package parse

import (
	"bufio"
	"errors"
	"go_redis/interface/resp"
	"go_redis/resp/reply"
	"io"
	"strconv"
	"strings"
)

type PayLoad struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	// 是多行还是单行解析
	readingMultiLine bool
	// 参数的数量
	exceptedArgsCount int
	// 消息类型
	msgType byte
	// 具体的数据
	args [][]byte
	// 数据块的长度
	bulkLen int64
}

// 解析是否完成
func (r *readState) isFinish() bool {
	// 参数全部解析完成
	return r.exceptedArgsCount > 0 && r.exceptedArgsCount == len(r.args)
}

// 异步解析

func ParseStream(reader io.Reader) <-chan *PayLoad {
	ch := make(chan *PayLoad)
	// 异步执行解析过程，通过管道来控制解析完成
	go parse0(reader, ch)
	return ch
}
func parse0(reader io.Reader, ch chan<- *PayLoad) {

}

// 读取一行
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 {
		// 1. 数据块为0，前面没有$符号指定的长度，读取到\r\n即可
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			// io错误
			return nil, true, err
		}
		// 判断msg结尾是否有\r
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			// 协议错误
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else {
		// 2. 前面出现过$
		msg = make([]byte, state.bulkLen+2)
		// 读取指定数量的字节
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			// io错误
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-1] != '\n' || msg[len(msg)-2] != '\r' {
			// 协议错误
			return nil, false, errors.New("protocol error: " + string(msg))
		}

	}
	return msg, false, nil

}

// 解析头指令(开头为*时执行)

func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine int64
	expectedLine, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if expectedLine == 0 {
		state.exceptedArgsCount = int(expectedLine)
		return nil
	} else if expectedLine > 0 {
		// 初始化state
		state.exceptedArgsCount = int(expectedLine)
		state.readingMultiLine = true
		state.args = make([][]byte, state.exceptedArgsCount)
		state.msgType = msg[0]
	} else {
		return errors.New("protocol error: " + string(msg))
	}
	return nil
}

// 解析一些信号数据
func parseSingle(msg []byte) (resp.Reply, error) {
	// 去除最后两位
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.MakeStatusReply(str[1:])
	case '-':
		result = reply.MakeStandErrReply(str[1:])
	case ':':
		val, err := strconv.Atoi(str[1:])
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		result = reply.MakeIntReply(int64(val))
	}
	return result, nil
}
