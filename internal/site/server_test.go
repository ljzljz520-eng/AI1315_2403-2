package site

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/pottery-experience/internal/booking"
)

func TestTwelvePagesAreAvailable(t *testing.T) {
	service := booking.NewUpdateService(booking.NewMemoryStore(booking.SeedRecords()), nil)
	handler := NewServer(service, "../../web/public").Handler()
	pages := []string{
		"/",
		"/courses.html",
		"/works.html",
		"/glazes.html",
		"/experience.html",
		"/experience-wheel.html",
		"/experience-cup.html",
		"/experience-plate.html",
		"/booking.html",
		"/studio.html",
		"/events.html",
		"/contact.html",
	}
	for _, page := range pages {
		request := httptest.NewRequest(http.MethodGet, page, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Errorf("%s returned %d", page, response.Code)
		}
	}
}
