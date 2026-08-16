package booking

func SeedRecords() []BookingRecord {
	return []BookingRecord{
		{
			ID:          "陶艺预约-2026-001",
			GuestName:   "林晓",
			Phone:       "13800001234",
			VisitDate:   "2026-09-12",
			Slot:        "14:00-16:00",
			ProjectID:   "wheel",
			ProjectName: "拉坯入门",
			Confirmations: map[string]string{
				"预约信息": "已确认",
			},
			Version: 1,
		},
	}
}
