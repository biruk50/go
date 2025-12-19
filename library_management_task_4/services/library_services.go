package services

import (
	"errors"
	"fmt"
	"library_management_task_4/models"

	"sync"
	"time"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	ReserveBook(bookID int, memberID int) error
}

// Library implements LibraryManager and supports concurrent reservations.
type Library struct {
	mu            sync.Mutex
	rwmu          sync.RWMutex
	Books         map[int]models.Book
	Members       map[int]models.Member
	reservations  map[int]int // bookID -> memberID (reserved)
	reserveTimers map[int]*time.Timer
}

// NewLibrary constructs an initialized library.
func NewLibrary() *Library {
	return &Library{
		Books:         make(map[int]models.Book),
		Members:       make(map[int]models.Member),
		reservations:  make(map[int]int),
		reserveTimers: make(map[int]*time.Timer),
	}
}

// AddBook adds or replaces a book in the library.
func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Books[book.Id] = book
}

// RemoveBook removes a book if it exists.
func (l *Library) RemoveBook(bookID int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if t, ok := l.reserveTimers[bookID]; ok && t != nil {
		t.Stop()
		delete(l.reserveTimers, bookID)
	}
	delete(l.reservations, bookID)
	delete(l.Books, bookID)
}

// It allows borrowing if the book is Available or Reserved for the same member.
func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}

	// If book is Reserved, ensure it's reserved for this member
	if book.Status == "Reserved" {
		reserver, reserved := l.reservations[bookID]
		if !reserved || reserver != memberID {
			return errors.New("book reserved by another member")
		}
	}

	book.Status = "Borrowed"
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Books[bookID] = book
	l.Members[memberID] = member

	// If there is a reservation/timer, cancel it and clear reservation
	if t, ok := l.reserveTimers[bookID]; ok && t != nil {
		t.Stop()
		delete(l.reserveTimers, bookID)
	}
	delete(l.reservations, bookID)

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	if book.Status == "Available" {
		return errors.New("book already returned (available)")
	}

	// Remove from member borrowed list
	newBorrowed := make([]models.Book, 0, len(member.BorrowedBooks))
	found := false
	for _, b := range member.BorrowedBooks {
		if b.Id == bookID {
			found = true
			continue
		}
		newBorrowed = append(newBorrowed, b)
	}
	if !found {
		return errors.New("book not borrowed by this member")
	}

	member.BorrowedBooks = newBorrowed
	l.Members[memberID] = member

	book.Status = "Available"
	l.Books[bookID] = book

	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.rwmu.RLock()
	defer l.rwmu.RUnlock()

	var available []models.Book
	for _, book := range l.Books {
		if book.Status == "Available" {
			available = append(available, book)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.rwmu.RLock()
	defer l.rwmu.RUnlock()

	member, ok := l.Members[memberID]
	if !ok {
		return []models.Book{}
	}
	return append([]models.Book(nil), member.BorrowedBooks...)
}

// ReserveBook reserves a book for a member and starts asynchronous borrow processing.
// The reservation auto-cancels if not borrowed within 5 seconds.

func (l *Library) ReserveBook(bookID int, memberID int) error {
	l.mu.Lock()

	book, ok := l.Books[bookID]
	if !ok {
		l.mu.Unlock()
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		l.mu.Unlock()
		return errors.New("book already borrowed")
	}
	if book.Status == "Reserved" {
		l.mu.Unlock()
		return errors.New("book already reserved")
	}
	if _, ok := l.Members[memberID]; !ok {
		l.mu.Unlock()
		return errors.New("member not found")
	}

	// Mark reserved and store reservation
	book.Status = "Reserved"
	l.Books[bookID] = book
	l.reservations[bookID] = memberID

	l.mu.Unlock() // release before spawning any goroutine

	// Timer: auto-cancel after 5 seconds
	timer := time.AfterFunc(5*time.Second, func() {
		l.mu.Lock()
		defer l.mu.Unlock()

		b, exists := l.Books[bookID]
		if !exists {
			return
		}
		// If still reserved and not borrowed, cancel reservation
		if b.Status == "Reserved" {
			delete(l.reservations, bookID)
			b.Status = "Available"
			l.Books[bookID] = b
			fmt.Printf("[Timer] Reservation for book %d auto-cancelled (member %d)\n", bookID, memberID)
		}
	})

	// Store timer safely
	l.mu.Lock()
	l.reserveTimers[bookID] = timer
	l.mu.Unlock()

	// Asynchronous borrow simulation
	go func() {
		time.Sleep(2 * time.Second)
		if err := l.BorrowBook(bookID, memberID); err != nil {
			fmt.Printf("[AsyncBorrow] Failed for book %d by member %d: %v\n", bookID, memberID, err)
			return
		}
		fmt.Printf("[AsyncBorrow] Book %d successfully borrowed by member %d\n", bookID, memberID)
	}()

	return nil
}
