package resp

// 回复接口
type Reply interface {
	ToBytes() []byte
}
