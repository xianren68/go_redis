package dict

import "sync"

type SyncDict struct {
	m sync.Map
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
	//TODO implement me
	panic("implement me")
}

func (s *SyncDict) Keys() []string {
	//TODO implement me
	panic("implement me")
}

func (s *SyncDict) RandomKeys(i int) []string {
	//TODO implement me
	panic("implement me")
}

func (s *SyncDict) RandomDistinctKeys(i int) []string {
	//TODO implement me
	panic("implement me")
}

func (s *SyncDict) Clear() {
	//TODO implement me
	panic("implement me")
}
