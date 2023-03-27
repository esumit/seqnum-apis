package seqnum

import (
	"context"
	"sync"
	"time"
)

const (
	// 41 bits for the timestamp (gives us 69 years of unique IDs)
	timeBits = 41

	// 22 bits for the counter
	counterBits = 22

	// Counter mask (all bits set to 1)
	counterMask = (1 << counterBits) - 1
)

type seqnumApiManager struct {
	mutex    sync.Mutex
	lastTime int64
	counter  int64
}

func NewSeqnumManager() *seqnumApiManager {
	return &seqnumApiManager{
		mutex:    sync.Mutex{},
		lastTime: 0,
		counter:  0,
	}
}

func (s *seqnumApiManager) GenerateSeqNum(ctx context.Context) (*SeqnumRs, error) {
	// Ensure thread-safety using a mutex
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get the current timestamp
	now := time.Now().UnixNano() / int64(time.Millisecond)

	// If the timestamp is the same as the last request, increment the counter
	if now == s.lastTime {
		s.counter = (s.counter + 1) & counterMask
	} else {
		// If the timestamp is different, reset the counter
		s.counter = 0
		s.lastTime = now
	}

	// Generate the unique sequence number
	seqNum := (s.lastTime << counterBits) | s.counter

	rs := SeqnumRs{
		SeqNum: seqNum,
	}
	return &rs, nil
}
