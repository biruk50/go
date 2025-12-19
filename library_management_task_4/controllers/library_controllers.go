package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"library_management_task_4/concurrency"
	"library_management_task_4/models"
	"library_management_task_4/services"
)

var print = fmt.Println

func RunLibraryConsole() {
	library := services.NewLibrary()

	library.Members[1] = models.Member{Id: 1, Name: "Abebe"}
	library.Members[2] = models.Member{Id: 2, Name: "Bob"}
	library.Members[3] = models.Member{Id: 3, Name: "Charlie"}

	library.AddBook(models.Book{Id: 1, Title: "alice in the wonderland", Author: "K. Cox"})
	library.AddBook(models.Book{Id: 2, Title: "game of thrones", Author: "A. Donovan"})

	// Start the reservation worker
	requestCh := concurrency.StartReservationWorker(library)

	reader := bufio.NewReader(os.Stdin)

	for {
		print("\n===== Library Management System =====")
		print("1. Add Book")
		print("2. Remove Book")
		print("3. Borrow Book (direct)")
		print("4. Return Book")
		print("5. List Available Books")
		print("6. List Borrowed Books by Member")
		print("7. Reserve Book (async via worker)")
		print("8. Simulate Concurrent Reservations")
		print("9. Exit")

		fmt.Print("Enter choice: ")
		input, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(input))

		switch choice {
		case 1:

			id := len(library.Books) + 1
			fmt.Print("Enter Title: ")
			title, _ := reader.ReadString('\n')
			fmt.Print("Enter Author: ")
			author, _ := reader.ReadString('\n')

			library.AddBook(models.Book{
				Id:     id,
				Title:  strings.TrimSpace(title),
				Author: strings.TrimSpace(author),
				Status: "Available",
			})

			print("Book added.")
		case 2:
			fmt.Print("Enter Book Id to remove: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			library.RemoveBook(id)
			print("Book removed.")
		case 3:
			fmt.Print("Enter Book Id: ")
			idStr, _ := reader.ReadString('\n')
			bookId, _ := strconv.Atoi(strings.TrimSpace(idStr))
			fmt.Print("Enter Member Id: ")
			memStr, _ := reader.ReadString('\n')
			memberId, _ := strconv.Atoi(strings.TrimSpace(memStr))

			err := library.BorrowBook(bookId, memberId)
			if err != nil {
				print("Error:", err)
			} else {
				print("Book borrowed.")
			}
		case 4:
			fmt.Print("Enter Book Id: ")
			idStr, _ := reader.ReadString('\n')
			bookId, _ := strconv.Atoi(strings.TrimSpace(idStr))
			fmt.Print("Enter Member Id: ")
			memStr, _ := reader.ReadString('\n')
			memberId, _ := strconv.Atoi(strings.TrimSpace(memStr))

			err := library.ReturnBook(bookId, memberId)
			if err != nil {
				print("Error:", err)
			} else {
				print("Book returned.")
			}
		case 5:
			books := library.ListAvailableBooks()
			print("Available Books:")
			for _, b := range books {
				fmt.Printf("Id: %d, Title: %s, Author: %s\n", b.Id, b.Title, b.Author)
			}
		case 6:
			fmt.Print("Enter Member Id: ")
			memStr, _ := reader.ReadString('\n')
			memberId, _ := strconv.Atoi(strings.TrimSpace(memStr))
			books := library.ListBorrowedBooks(memberId)
			fmt.Printf("Books borrowed by Member %d:\n", memberId)
			for _, b := range books {
				fmt.Printf("Id: %d, Title: %s\n", b.Id, b.Title)
			}
		case 7:
			fmt.Print("Enter Book Id to reserve: ")
			idStr, _ := reader.ReadString('\n')
			bookId, _ := strconv.Atoi(strings.TrimSpace(idStr))
			fmt.Print("Enter Member Id: ")
			memStr, _ := reader.ReadString('\n')
			memberId, _ := strconv.Atoi(strings.TrimSpace(memStr))

			resp := make(chan error)
			requestCh <- concurrency.ReservationRequest{BookID: bookId, MemberID: memberId, Resp: resp}
			err := <-resp
			if err != nil {
				print("Reserve Error:", err)
			} else {
				print("Reserved successfully (async borrow will be attempted).")
			}

		case 8:
			// Simulate multiple members trying to reserve the same book concurrently
			fmt.Print("Enter Book Id to simulate concurrent reservations: ")
			idStr, _ := reader.ReadString('\n')
			bookId, _ := strconv.Atoi(strings.TrimSpace(idStr))
			fmt.Print("Enter comma-separated Member Ids (e.g., 1,2,3): ")
			memLine, _ := reader.ReadString('\n')
			parts := strings.Split(strings.TrimSpace(memLine), ",")

			respChans := make([]chan error, 0, len(parts))
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}

				memberId, err := strconv.Atoi(p)
				if err != nil {
					fmt.Printf("invalid member id: %s\n", p)
					continue
				}
				resp := make(chan error)
				respChans = append(respChans, resp)
				requestCh <- concurrency.ReservationRequest{BookID: bookId, MemberID: memberId, Resp: resp}
				// tiny stagger to better simulate real concurrency (optional)
				time.Sleep(10 * time.Millisecond)
			}
			// Collect responses
			for i, ch := range respChans {
				err := <-ch
				if err != nil {
					fmt.Printf("Request %d error: %v\n", i+1, err)
				} else {
					fmt.Printf("Request %d succeeded (reservation accepted, async borrow attempted)\n", i+1)
				}
			}
			// Allow some time to see async borrow results & timers in console
			print("Waiting 5 seconds to allow async borrows / auto-cancellations to appear...")
			time.Sleep(5 * time.Second)
		case 9:
			print("Exiting system. Goodbye!")
			// close worker channel to stop worker goroutine cleanly
			close(requestCh)
			return
		default:
			print("Invalid choice")
		}
	}
}
