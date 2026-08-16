package booking

import "errors"

var (
	ErrRecordNotFound = errors.New("booking record not found")
	ErrInvalidChange  = errors.New("invalid operator change")
)

type BookingRecord struct {
	ID            string            `json:"id"`
	GuestName     string            `json:"guestName"`
	Phone         string            `json:"phone"`
	VisitDate     string            `json:"visitDate"`
	Slot          string            `json:"slot"`
	ProjectID     string            `json:"projectId"`
	ProjectName   string            `json:"projectName"`
	Confirmations map[string]string `json:"confirmations"`
	Version       int               `json:"version"`
}

type OperatorChange struct {
	Operator string `json:"operator"`
	Field    string `json:"field"`
	Value    string `json:"value"`
}

type Summary struct {
	ID            string            `json:"id"`
	GuestName     string            `json:"guestName"`
	VisitDate     string            `json:"visitDate"`
	Slot          string            `json:"slot"`
	ProjectName   string            `json:"projectName"`
	Confirmations map[string]string `json:"confirmations"`
	Version       int               `json:"version"`
}

func (r BookingRecord) Summary() Summary {
	confirmations := make(map[string]string, len(r.Confirmations))
	for key, value := range r.Confirmations {
		confirmations[key] = value
	}
	return Summary{
		ID:            r.ID,
		GuestName:     r.GuestName,
		VisitDate:     r.VisitDate,
		Slot:          r.Slot,
		ProjectName:   r.ProjectName,
		Confirmations: confirmations,
		Version:       r.Version,
	}
}
