package booking

import (
	"strings"
	"sync/atomic"
)

type UpdateService struct {
	store *MemoryStore
	gate  *UpdateGate
}

func NewUpdateService(store *MemoryStore, gate *UpdateGate) *UpdateService {
	return &UpdateService{store: store, gate: gate}
}

func (s *UpdateService) Confirm(id string, change OperatorChange) (BookingRecord, error) {
	if strings.TrimSpace(change.Operator) == "" || strings.TrimSpace(change.Field) == "" || strings.TrimSpace(change.Value) == "" {
		return BookingRecord{}, ErrInvalidChange
	}
	if s.gate != nil {
		s.gate.BeforeWrite()
	}
	// The read-modify-write happens atomically inside the store so that two
	// concurrent Confirm calls (even when the gate makes them arrive together)
	// merge instead of overwriting each other's confirmations.
	return s.store.Update(id, func(record BookingRecord) BookingRecord {
		if record.Confirmations == nil {
			record.Confirmations = make(map[string]string)
		}
		record.Confirmations[change.Field] = change.Value
		record.Version++
		return record
	})
}

func (s *UpdateService) GetSummary(id string) (Summary, error) {
	record, err := s.store.Get(id)
	if err != nil {
		return Summary{}, err
	}
	return record.Summary(), nil
}

type UpdateGate struct {
	expected int32
	arrived  int32
	ready    chan struct{}
	release  chan struct{}
}

func NewUpdateGate(expected int) *UpdateGate {
	return &UpdateGate{
		expected: int32(expected),
		ready:    make(chan struct{}),
		release:  make(chan struct{}),
	}
}

func (g *UpdateGate) BeforeWrite() {
	if atomic.AddInt32(&g.arrived, 1) == g.expected {
		close(g.ready)
	}
	<-g.release
}

func (g *UpdateGate) WaitUntilReady() {
	<-g.ready
}

func (g *UpdateGate) Release() {
	close(g.release)
}
