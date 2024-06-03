package parser

import "sync"

type subscriber interface {
	Subscribe(address string) bool
	Observer
}

// Observer ...
type Observer interface {
	GetTransactions(address string) []Transaction
	UpdateTransactions(address string, txs []Transaction)
	GetSubscribers() map[string][]Transaction
}

// MemStorage ...
type MemStorage struct {
	lock      sync.Mutex
	addresses map[string][]Transaction
}

// Subscribe implements the subscriber interface
func (m *MemStorage) Subscribe(address string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.addresses[address]; ok {
		return false
	}
	m.addresses[address] = []Transaction{}
	return true
}

// UpdateTransactions implements the Observer interface
func (m *MemStorage) UpdateTransactions(address string, txs []Transaction) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.addresses[address]; !ok {
		return
	}
	// overwrite, this is a simple implementation
	m.addresses[address] = txs
}

// GetTransactions implements the Observer interface
func (m *MemStorage) GetTransactions(addr string) []Transaction {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.addresses[addr]; !ok {
		return nil
	}
	return m.addresses[addr]
}

// GetSubscribers implements the Observer interface
func (m *MemStorage) GetSubscribers() map[string][]Transaction {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.addresses
}
