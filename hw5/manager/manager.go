package manager

import (
	"errors"
	"fmt"
	"sync"
)

type Shard struct {
	Address string
	Number  int
}

// Struct to balance between master and slave
type balancer struct {
	sync.Mutex
	balance int64
}

type Manager struct {
	size int
	ss   *sync.Map
	b    *balancer
}

var (
	ErrorShardNotFound = errors.New("shard not found")
)

func NewManager(size int) *Manager {
	return &Manager{
		size: size,
		ss:   &sync.Map{},
		b:    &balancer{balance: 0},
	}
}
func (m *Manager) Add(s *Shard) {
	m.ss.Store(s.Number, s)
}
func (m *Manager) ShardById(entityId int) (*Shard, error) {
	if entityId < 0 {
		return nil, ErrorShardNotFound
	}
	n := entityId % m.size // TODO: think about devision to shards
	m.b.Lock()
	if m.b.balance > 0 {
		m.b.balance--
		n += 10
	} else {
		m.b.balance++
	}
	m.b.Unlock()
	if s, ok := m.ss.Load(n); ok {
		fmt.Println(s)
		return s.(*Shard), nil
	}
	return nil, ErrorShardNotFound
}
