package booking

import "sync"

type MemoryStore struct {
	mu      sync.RWMutex
	records map[string]BookingRecord
}

func NewMemoryStore(records []BookingRecord) *MemoryStore {
	store := &MemoryStore{records: make(map[string]BookingRecord, len(records))}
	for _, record := range records {
		store.records[record.ID] = cloneRecord(record)
	}
	return store
}

func (s *MemoryStore) Get(id string) (BookingRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.records[id]
	if !ok {
		return BookingRecord{}, ErrRecordNotFound
	}
	return cloneRecord(record), nil
}

func (s *MemoryStore) Replace(record BookingRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[record.ID]; !ok {
		return ErrRecordNotFound
	}
	s.records[record.ID] = cloneRecord(record)
	return nil
}

// Update performs an atomic read-modify-write: apply receives an isolated
// clone of the current record and returns the updated record, all while the
// store's write lock is held. This prevents two concurrent Confirm calls from
// reading the same stale version and overwriting each other's changes (a lost
// update). The gate may still force both callers to arrive together, but the
// lock serializes their mutations so every valid change is merged and persisted.
func (s *MemoryStore) Update(id string, apply func(BookingRecord) BookingRecord) (BookingRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.records[id]
	if !ok {
		return BookingRecord{}, ErrRecordNotFound
	}
	updated := apply(cloneRecord(record))
	s.records[id] = cloneRecord(updated)
	return cloneRecord(updated), nil
}

func cloneRecord(record BookingRecord) BookingRecord {
	copyOf := record
	copyOf.Confirmations = make(map[string]string, len(record.Confirmations))
	for key, value := range record.Confirmations {
		copyOf.Confirmations[key] = value
	}
	return copyOf
}
