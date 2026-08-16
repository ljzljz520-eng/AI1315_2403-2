package booking

import (
	"sync"
	"testing"
)

func TestSingleOperatorConfirmation(t *testing.T) {
	service := NewUpdateService(NewMemoryStore(SeedRecords()), nil)
	_, err := service.Confirm("陶艺预约-2026-001", OperatorChange{
		Operator: "周老师",
		Field:    "泥料",
		Value:    "白陶泥",
	})
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}
	summary, err := service.GetSummary("陶艺预约-2026-001")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.Confirmations["泥料"] != "白陶泥" {
		t.Fatalf("unexpected confirmation: %#v", summary.Confirmations)
	}
}

func TestConcurrentOperatorConfirmationsAreMerged(t *testing.T) {
	gate := NewUpdateGate(2)
	service := NewUpdateService(NewMemoryStore(SeedRecords()), gate)
	changes := []OperatorChange{
		{Operator: "周老师", Field: "泥料", Value: "白陶泥"},
		{Operator: "许老师", Field: "釉色", Value: "青瓷釉"},
	}
	var wg sync.WaitGroup
	for _, change := range changes {
		change := change
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := service.Confirm("陶艺预约-2026-001", change); err != nil {
				t.Errorf("confirm: %v", err)
			}
		}()
	}
	gate.WaitUntilReady()
	gate.Release()
	wg.Wait()

	summary, err := service.GetSummary("陶艺预约-2026-001")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.Confirmations["泥料"] != "白陶泥" || summary.Confirmations["釉色"] != "青瓷釉" {
		t.Fatalf("summary lost a confirmation: %#v", summary.Confirmations)
	}
}
