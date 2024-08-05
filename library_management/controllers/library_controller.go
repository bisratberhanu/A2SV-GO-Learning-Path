package controllers

import (
    "bufio"
    "fmt"
    "library_management/models"
    "library_management/services"
    "os"
    "strconv"
)

var library = services.NewLibrary()

func StartLibraryController() {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Println("1. Add a new book")
        fmt.Println("2. Remove an existing book")
        fmt.Println("3. Borrow a book")
        fmt.Println("4. Return a book")
        fmt.Println("5. List all available books")
        fmt.Println("6. List all borrowed books by a member")
        fmt.Println("7. Add a new member")
        fmt.Println("8. Exit")
        fmt.Print("Enter your choice: ")

        scanner.Scan()
        choice := scanner.Text()

        switch choice {
        case "1":
            addBook(scanner)
        case "2":
            removeBook(scanner)
        case "3":
            borrowBook(scanner)
        case "4":
            returnBook(scanner)
        case "5":
            listAvailableBooks()
        case "6":
            listBorrowedBooks(scanner)
        case "7":
            addMember(scanner)
        case "8":
            fmt.Println("Exiting...")
            return
        default:
            fmt.Println("Invalid choice, please try again.")
        }
    }
}

func addBook(scanner *bufio.Scanner) {
    fmt.Print("Enter book ID: ")
    scanner.Scan()
    bookID, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid book ID")
        return
    }

    fmt.Print("Enter book title: ")
    scanner.Scan()
    bookTitle := scanner.Text()

    fmt.Print("Enter book author: ")
    scanner.Scan()
    bookAuthor := scanner.Text()

    book := models.Book{
        ID:     bookID,
        Title:  bookTitle,
        Author: bookAuthor,
        Status: "Available",
    }
    library.AddBook(book)
    fmt.Println("Book added successfully.")
}

func removeBook(scanner *bufio.Scanner) {
    fmt.Print("Enter book ID to remove: ")
    scanner.Scan()
    bookID, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid book ID")
        return
    }

    library.RemoveBook(bookID)
    fmt.Println("Book removed successfully.")
}

func borrowBook(scanner *bufio.Scanner) {
    fmt.Print("Enter book ID to borrow: ")
    scanner.Scan()
    bookID, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid book ID")
        return
    }

    fmt.Print("Enter member ID: ")
    scanner.Scan()
    memberID, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid member ID")
        return
    }

    err = library.BorrowBook(bookID, memberID)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Book borrowed successfully.")
    }
}

func returnBook(scanner *bufio.Scanner) {
    fmt.Print("Enter book ID to return: ")
    scanner.Scan()
    bookID, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid book ID")
        return
    }

    fmt.Print("Enter member ID: ")
    scanner.Scan()
    memberID, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid member ID")
        return
    }

    err = library.ReturnBook(bookID, memberID)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Book returned successfully.")
    }
}

func listAvailableBooks() {
    books := library.ListAvailableBooks()
    if len(books) == 0 {
        fmt.Println("No available books.")
    } else {
        fmt.Println("Available books:")
        for _, book := range books {
            fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
        }
    }
}

func listBorrowedBooks(scanner *bufio.Scanner) {
    fmt.Print("Enter member ID: ")
    scanner.Scan()
    memberID, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid member ID")
        return
    }

    books := library.ListBorrowedBooks(memberID)
    if len(books) == 0 {
        fmt.Println("No borrowed books for this member.")
    } else {
        fmt.Println("Borrowed books:")
        for _, book := range books {
            fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
        }
    }
}

func addMember(scanner *bufio.Scanner) {
    fmt.Print("Enter member ID: ")
    scanner.Scan()
    memberID, err := strconv.Atoi(scanner.Text())
    if err != nil {
        fmt.Println("Invalid member ID")
        return
    }

    fmt.Print("Enter member name: ")
    scanner.Scan()
    memberName := scanner.Text()

    member := models.Member{
        ID:           memberID,
        Name:         memberName,
        BorrowedBooks: []models.Book{},
    }
    library.AddMember(member)
    fmt.Println("Member added successfully.")
}
