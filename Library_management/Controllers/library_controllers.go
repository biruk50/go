package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

var print = fmt.Println

func StartConsole(library *services.Library) {
	reader := bufio.NewReader(os.Stdin)

	for {
		print("\nLibrary Management System")
		print("1. Add Book")
		print("2. Remove Book")
		print("3. Borrow Book")
		print("4. Return Book")
		print("5. List Available Books")
		print("6. List Borrowed Books")
		print("7. Exit")
		fmt.Print("Enter choice: ")

		choiceStr, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(choiceStr))

		switch choice {
		case 1:
			id := len(library.Books) + 1
			fmt.Print("Enter Title: ")
			title := readLine(reader)
			fmt.Print("Enter Author: ")
			author := readLine(reader)

			library.AddBook(models.Book{Id: id, Title: title, Author: author, Status: "Available"})
			print("Book added successfully!")

		case 2:
			fmt.Print("Enter Book ID to remove: ")
			id, _ := strconv.Atoi(strings.TrimSpace(readLine(reader)))
			library.RemoveBook(id)
			print("Book removed successfully!")

		case 3:
			fmt.Print("Enter Member ID: ")
			memberID, _ := strconv.Atoi(strings.TrimSpace(readLine(reader)))
			fmt.Print("Enter Book ID: ")
			bookID, _ := strconv.Atoi(strings.TrimSpace(readLine(reader)))

			if _, ok := library.Members[memberID]; !ok {
				fmt.Print("Enter Member Name: ")
				name := readLine(reader)
				library.Members[memberID] = models.Member{Id: memberID, Name: name}
			}

			err := library.BorrowBook(bookID, memberID)
			if err != nil {
				print("Error:", err)
			} else {
				print("Book borrowed successfully!")
			}

		case 4:
			fmt.Print("Enter Member ID: ")
			memberID, _ := strconv.Atoi(strings.TrimSpace(readLine(reader)))
			fmt.Print("Enter Book ID: ")
			bookID, _ := strconv.Atoi(strings.TrimSpace(readLine(reader)))

			err := library.ReturnBook(bookID, memberID)
			if err != nil {
				print("Error:", err)
			} else {
				print("Book returned successfully!")
			}

		case 5:
			print("\nAvailable Books:")
			for _, b := range library.ListAvailableBooks() {
				fmt.Printf("[%d] %s by %s\n", b.Id, b.Title, b.Author)
			}

		case 6:
			fmt.Print("Enter Member ID: ")
			memberID, _ := strconv.Atoi(strings.TrimSpace(readLine(reader)))
			print("\nBorrowed Books:")
			for _, b := range library.ListBorrowedBooks(memberID) {
				fmt.Printf("[%d] %s by %s\n", b.Id, b.Title, b.Author)
			}

		case 7:
			print("Goodbye!")
			return
		default:
			print("Invalid choice, try again.")
		}
	}
}

func readLine(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
