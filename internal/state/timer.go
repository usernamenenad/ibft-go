package state

import "time"

type RoundChangeTimer struct {
	IsRunning     bool
	ExpiryTime    time.Time
	ExpiryHandler func()
	ticker        *time.Ticker
	stopCh        chan struct{}
}

func NewRoundChangeTimer() *RoundChangeTimer {
	return &RoundChangeTimer{
		IsRunning: false,
		stopCh:    make(chan struct{}),
	}
}

func (t *RoundChangeTimer) Start(timeout time.Duration, expiryHandler func()) {
	t.Stop()
	t.IsRunning = true
	t.ExpiryTime = time.Now().Add(timeout)
	t.ExpiryHandler = expiryHandler
	t.ticker = time.NewTicker(timeout)

	go func() {
		select {
		case <-t.ticker.C:
			if t.ExpiryHandler != nil {
				t.ExpiryHandler()
			}
		case <-t.stopCh:
			return
		}
	}()
}

func (t *RoundChangeTimer) Stop() {
	if t.IsRunning && t.ticker != nil {
		t.ticker.Stop()
		t.IsRunning = false
		select {
		case t.stopCh <- struct{}{}:
		default:
		}
	}
}

func (t *RoundChangeTimer) calculateTimeout(round uint64) time.Duration {
	return 2000
}
