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

// Update applies a change to the current record while holding the store lock,
// so concurrent changes cannot overwrite fields committed by another update.
func (s *MemoryStore) Update(id string, update func(*BookingRecord) error) (BookingRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.records[id]
	if !ok {
		return BookingRecord{}, ErrRecordNotFound
	}
	record = cloneRecord(record)
	if err := update(&record); err != nil {
		return BookingRecord{}, err
	}
	s.records[id] = cloneRecord(record)
	return cloneRecord(record), nil
}

func cloneRecord(record BookingRecord) BookingRecord {
	copyOf := record
	copyOf.Confirmations = make(map[string]string, len(record.Confirmations))
	for key, value := range record.Confirmations {
		copyOf.Confirmations[key] = value
	}
	return copyOf
}
