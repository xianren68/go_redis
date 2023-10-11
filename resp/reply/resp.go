package reply

import (
	"bytes"
	"go_redis/interface/resp"
	"strconv"
)

// 开始与结束的标识符

var CRLF = "\r\n"

/* ————— Bulk ————— */

type BulkReply struct {
	Arg []byte
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		arg,
	}
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkBytes
	}

	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

/* ———— Multi ———— */

type MultiBulkReply struct {
	Args [][]byte
}

func MakeMultiBulkReply(arg [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		arg,
	}
}

func (m *MultiBulkReply) ToBytes() []byte {
	l := len(m.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(l) + CRLF)
	for _, val := range m.Args {
		if len(val) == 0 {
			buf.WriteString("$-1" + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(val)) + CRLF + string(val) + CRLF)
		}

	}
	return buf.Bytes()

}

/* ———— Status ———— */

type StatusReply struct {
	Status string
}

func (s *StatusReply) ToBytes() []byte {
	return []byte(s.Status)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

/* ———— Int ———— */

type IntReply struct {
	Code int64
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		code,
	}
}

/* ———— Error ———— */

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

type StandErrReply struct {
	Status string
}

func MakeStandErrReply(status string) *StatusReply {
	return &StatusReply{
		status,
	}
}

func IsErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}

func (s *StandErrReply) Error() string {
	return s.Status
}

func (s *StandErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}
