package reply

/* 一些固定回复 */
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func MakePongReply() *PongReply {
	return &PongReply{}
}
func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

func MakeOkReply() *OkReply {
	return &OkReply{}
}
func (ok *OkReply) ToBytes() []byte {
	return okBytes
}

type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func MakeNullBulkBytes() *NullBulkReply {
	return &NullBulkReply{}
}
func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}
func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

type NoReply struct {
}

var NoBytes = []byte("")

func MakeNoReply() *NoReply {
	return &NoReply{}
}
func (n *NoReply) ToBytes() []byte {
	return NoBytes
}
