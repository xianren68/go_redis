package resp

// 未知错误
type UnKnowErrReply struct {
}

var unKnowErrBytes = []byte("-Err unkonw\r\n")

func (u *UnKnowErrReply) Error() string {
	return string(unKnowErrBytes)
}

func (u *UnKnowErrReply) ToBytes() []byte {
	return unKnowErrBytes
}

// 参数错误
type ArgNumErrReply struct {
	// 错误的指令
	Cmd string
}

func (a *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + a.Cmd + "' command"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + a.Cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		cmd,
	}
}

// 句法错误
type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

func MakeSyntaxErrReply() *SyntaxErrReply {
	return &SyntaxErrReply{}
}

func (r *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

// 类型错误
type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

// 协议错误（不符合规范）
type ProtocolErrReply struct {
	Msg string
}

func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "'\r\n")
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error '" + r.Msg + "' command"
}
