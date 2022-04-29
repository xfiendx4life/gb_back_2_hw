package manager

import (
	"errors"
	"log"
	"sync"
)

type Shard struct {
	Role    string
	Address string
	Number  int
}
type Manager struct {
	size int
	ss   *sync.Map
}

var (
	ErrorShardNotFound = errors.New("shard not found")
)

func NewManager(size int) *Manager {
	return &Manager{
		size: size,
		ss:   &sync.Map{},
	}
}
func (m *Manager) Add(s *Shard) {
	m.ss.Store(s.Number, s)
}
func (m *Manager) ShardById(entityId int, master bool) (*Shard, error) {
	if entityId < 0 {
		return nil, ErrorShardNotFound
	}
	n := entityId % m.size // TODO: think about devision to shards
	if !master {
		n += 10
	}
	if s, ok := m.ss.Load(n); ok {
		sh := s.(*Shard)
		log.Printf("operation on shard #%d role: %s\n", sh.Number, sh.Role)
		return s.(*Shard), nil
	}
	return nil, ErrorShardNotFound
}
