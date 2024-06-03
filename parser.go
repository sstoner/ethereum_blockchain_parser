package parser

import (
	"context"
	"math/rand"
	"time"
)

// Explorer implements the Parser interface
type Explorer struct {
	// service is the blockchain service
	service *Ethereum
	// subscriber is the state of subscribers
	subscriber subscriber
}

// GetCurrentBlock implements the Parser interface
func (e *Explorer) GetCurrentBlock() int {
	curBlock, err := e.service.GetCurrentBlock()
	if err != nil {
		// TODO expose error
		return -1
	}
	return curBlock
}

// Subscribe implements the Parser interface
func (e *Explorer) Subscribe(address string) bool {
	if !e.subscriber.Subscribe(address) {
		return false
	}

	txs, err := e.service.GetTransactions(address)
	if err != nil {
		// TODO expose error
		return false
	}

	e.subscriber.UpdateTransactions(address, txs)
	return true
}

// GetTransactions implements the Parser interface
func (e *Explorer) GetTransactions(address string) []Transaction {
	// return from cached
	return e.subscriber.GetTransactions(address)
}

// Observer should be run as a goroutine
func (e *Explorer) Observer(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			subs := e.subscriber.GetSubscribers()
			for addr := range subs {
				txs, err := e.service.GetTransactions(addr)
				if err != nil {
					// TODO log error
					continue
				}

				e.subscriber.UpdateTransactions(addr, txs)
			}
		}
	}
}

// NewExplorer returns a new explorer
func NewExplorer(subscriber subscriber) *Explorer {
	ethService := NewEthereum(ethEndpoint, httpTimeout)
	return &Explorer{
		service:    ethService,
		subscriber: subscriber,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	ethEndpoint = "https://cloudflare-eth.com"
	httpTimeout = 5 * time.Second
)
