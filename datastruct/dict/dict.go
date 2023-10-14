package dict

type Consumer func(any, any) bool
type Dict interface {
	Get(string) (any, bool)
	Len() int
	// Put 插入数据
	Put(string, any) int
	// PutIfAbsent 插入数据，如果不存在的话
	PutIfAbsent(string, any) int
	// PutIfExist 更新数据（只有数据存在时才插入）
	PutIfExist(string, any) int
	// Remove 删除数据
	Remove(string) int
	// ForEach 遍历数据
	ForEach(Consumer)
	Keys() []string
	// RandomKeys 随机返回指定数量的键
	RandomKeys(int) []string
	// RandomDistinctKeys 随机返回指定数量的键 不重复
	RandomDistinctKeys(int) []string
	Clear()
}
