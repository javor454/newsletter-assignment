package healthcheck

import (
	"context"
	"sync"
	"time"
)

type Status string

type HealthIndicator interface {
	Name() string
	Check(ctx context.Context) bool
}

type HealthMonitor struct {
	indicators []HealthIndicator
	status     map[string]bool
	mu         sync.RWMutex
}

func NewHealthMonitor(indicators ...HealthIndicator) *HealthMonitor {
	hm := &HealthMonitor{
		indicators: indicators,
		status:     make(map[string]bool),
	}
	for _, ind := range indicators {
		hm.status[ind.Name()] = true
	}

	hm.start(5 * time.Second)

	return hm
}

func (hm *HealthMonitor) start(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			hm.checkAll()
		}
	}()
}

func (hm *HealthMonitor) checkAll() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, ind := range hm.indicators {
		status := ind.Check(ctx)
		hm.mu.Lock()
		hm.status[ind.Name()] = status
		hm.mu.Unlock()
	}
}

func (hm *HealthMonitor) GetStatus() map[string]bool {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	status := make(map[string]bool)
	for k, v := range hm.status {
		status[k] = v
	}

	return status
}
