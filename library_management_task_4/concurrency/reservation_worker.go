package concurrency

import (
	"library_management_task_4/services"
)

// ReservationRequest is sent to the worker to reserve a book.
type ReservationRequest struct {
	BookID  int
	MemberID int
	Resp    chan error
}

// StartReservationWorker starts a goroutine that listens for reservation requests.
// It returns a channel where callers can send ReservationRequest items.
func StartReservationWorker(l *services.Library) chan ReservationRequest {
	ch := make(chan ReservationRequest)
	go func() {
		for req := range ch {
			// handle each request concurrently to simulate high load; library methods are concurrency-safe
			go func(r ReservationRequest) {
				err := l.ReserveBook(r.BookID, r.MemberID)
				r.Resp <- err
			}(req)
		}
	}()
	return ch
}
