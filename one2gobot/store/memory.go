package store

import (
	"greenisha/one2gobot/model"
)

type Session struct {
	ChatId      int64
	UserId      int64
	StationFrom model.Station
	StationTo   model.Station
	Mode        Mode
}

type Mode int

const (
	SelectFrom Mode = iota
	SelectTo
	Finished
)

type Store interface {
	Set(s Session)
	Get(u int64) (Session, bool)
	Delete(u int64)
}

type MemoryStore struct {
	data map[int64]Session
}

func (m *MemoryStore) Set(s Session) {
	m.data[s.UserId] = s

}
func (m *MemoryStore) Delete(u int64) {
	delete(m.data, u)

}
func (m *MemoryStore) Get(u int64) (Session, bool) {
	val, found := m.data[u]
	if !found {
		return Session{}, false
	}
	return val, true
}
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[int64]Session)}
}
