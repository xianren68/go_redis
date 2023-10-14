package dict

import (
	"sync"
)

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{
		m: sync.Map{},
	}
}
func (s *SyncDict) Get(a string) (any, bool) {
	return s.m.Load(a)
}

func (s *SyncDict) Len() int {
	length := 0
	s.m.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

func (s *SyncDict) Put(s2 string, a any) int {
	_, existed := s.Get(s2)
	if existed {
		return 0
	}
	s.m.Store(s2, a)
	return 1
}

func (s *SyncDict) PutIfAbsent(s2 string, a any) int {
	_, existed := s.Get(s2)
	if existed {
		return 0
	}
	s.m.Store(s2, a)
	return 1
}

func (s *SyncDict) PutIfExist(s2 string, a any) int {
	_, existed := s.Get(s2)
	if existed {
		s.m.Store(s2, a)
		return 1
	}
	return 0
}

func (s *SyncDict) Remove(s2 string) int {
	_, existed := s.Get(s2)
	if existed {
		s.m.Delete(s2)
		return 1
	}
	return 0
}

func (s *SyncDict) ForEach(consumer Consumer) {
	s.m.Range(consumer)
}

func (s *SyncDict) Keys() []string {
	l := s.Len()
	res := make([]string, 0, l)
	s.m.Range(func(key, value any) bool {
		res = append(res, key.(string))
		return true
	})
	return res
}

func (s *SyncDict) RandomKeys(limit int) []string {
	res := make([]string, limit)
	for i := 0; i < limit; i++ {
		s.m.Range(func(key, value any) bool {
			res[i] = key.(string)
			return false
		})
	}
	return res
}

func (s *SyncDict) RandomDistinctKeys(limit int) []string {
	res := make([]string, limit)
	i := 0
	s.m.Range(func(key, value any) bool {
		if i == limit {
			return false
		}
		res[i] = key.(string)
		i++
		return true
	})
	return res
}

func (s *SyncDict) Clear() {
	s.m = sync.Map{}
}
