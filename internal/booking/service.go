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
	record, err := s.store.Get(id)
	if err != nil {
		return BookingRecord{}, err
	}
	if s.gate != nil {
		s.gate.BeforeWrite()
	}
	if record.Confirmations == nil {
		record.Confirmations = make(map[string]string)
	}
	record.Confirmations[change.Field] = change.Value
	record.Version++
	if err := s.store.Replace(record); err != nil {
		return BookingRecord{}, err
	}
	return record, nil
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
