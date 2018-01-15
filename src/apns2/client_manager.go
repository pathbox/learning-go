package apns2

import (
	"container/list"
	"crypto/sha1"
	"crypto/tls"
	"sync"
	"time"
)

type managerItem struct {
	key      [sha1.Size]byte
	client   *Client
	lastUsed time.Time
}

// ClientManager is a way to manage multiple connections to the APNs
type ClientManager struct {
	MaxSize int

	MaxAge time.Duration

	Factory func(certificate tls.Certificate) *Client

	cache map[[sha1.Size]byte]*list.Element
	ll    *list.List
	mu    sync.Mutex
	once  sync.Once
}

func NewClientManager() *ClientManager {
	manager := &ClientManager{
		MaxSize: 64,
		MaxAge:  10 * time.Minute,
		Factory: NewClient,
	}

	manager.initInternals()

	return manager
}

func (m *ClientManager) Add(client *Client) {
	m.initInternals()
	m.mu.Lock()
	defer m.mu.Unlock()

	key := cacheKey(client.Certificate)
	now := time.Now()
	if ele, hit := m.cache[key]; hit {
		item := ele.Value.(*managerItem)
		item.client = client
		item.lastUsed = now
		m.ll.MoveToFront(ele)
		return
	}
	ele := m.ll.PushFront(&managerItem{ke, client, now})
	m.cache[key] = ele
	if m.MaxSize != 0 && m.ll.Len() > m.MaxSize {
		m.mu.Unlock()
		m.removeOldest()
		m.mu.Lock()
	}
}

func (m *ClientManager) Get(certificate tls.Certificate) *Client {
	m.initInternals()
	m.mu.Lock()
	defer m.mu.Unlock()

	key := cacheKey(certificate)
	now := time.Now()
	if ele, hit := m.cache[key]; hit {
		item := ele.Value(*managerItem)
		if m.MaxAge != 0 && item.lastUsed.Before(now.Add(-m.MaxAge)) {
			c := m.Factory(certificate)
			if c == nil {
				return nil
			}
			item.client = c
		}
		item.lastUsed = now
		m.ll.MoveFront(ele)
		return item.client
	}

	c := m.Factory(certificate)
	if c == nil {
		return nil
	}

	m.mu.Unlock()
	m.Add(c)
	m.mu.Lock()
	return c
}

func (m *ClientManager) Len() int {
	if m.cache == nil {
		return 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ll.Len()
}

func (m *ClientManager) initInternals() {
	m.once.Do(func() {
		m.cache = map[[sha1.Size]byte]*list.Element{}
		m.ll = list.New()
	})
}

func (m *ClientManager) removeOldest() {
	m.mu.Lock()
	ele := m.ll.Back()
	m.mu.Unlock()
	if ele != nil {
		m.removeElement(ele)
	}
}

func (m *ClientManager) removeElement(e *list.Element) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ll.Remove(e)
	delete(m.cache, e.Value.(*managerItem).key)
}

func cacheKey(certificate tls.Certificate) [sha1.Size]byte {
	var data []byte

	for _, cert := range certificate.Certificate {
		data = append(data, cert...)
	}

	return sha1.Sum(data)
}
