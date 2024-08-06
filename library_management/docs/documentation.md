# Library Management System

## Introduction
This document provides comprehensive information about the console-based Library Management System implemented in Go. The system demonstrates the use of structs, interfaces, methods, slices, and maps in Go. It allows for managing books and members in a library through various functionalities.

## Features Implemented
1. **Add Book**: Adds a new book to the library.
2. **Remove Book**: Removes an existing book from the library by its ID.
3. **Borrow Book**: Allows a member to borrow a book if it is available.
4. **Return Book**: Allows a member to return a borrowed book.
5. **List Available Books**: Lists all available books in the library.
6. **List Borrowed Books**: Lists all books borrowed by a specific member.
7. **Add Member**: Adds a new member to the library.

## Folder Structure
The project follows the below folder structure:

library_management/

├── main.go

├── controllers/

│ └── library_controller.go

├── models/

│ └── book.go

│ └── member.go

├── services/

│ └── library_service.go

├── docs/

│ └── documentation.md

└── go.mod


### main.go
The entry point of the application. It starts the library controller which handles user interactions.

```go
package main

import (
    "library_management/controllers"
)

func main() {
    controllers.StartLibraryController()
}

```

### controllers/library_controller
Handles console input and invokes the appropriate service methods. It provides functions for adding, removing, borrowing, returning books, listing available books, listing borrowed books by a member, and adding new members.

### models/book.go
Defines the Book struct with fields ID, Title, Author, and Status.
```go
package models

type Book struct {
    ID     int
    Title  string
    Author string
    Status string // can be "Available" or "Borrowed"
}
```
### models/member.go
Defines the Member struct with fields ID, Name, and BorrowedBooks.

```go
package models

type Member struct {
    ID           int
    Name         string
    BorrowedBooks []Book
}
```
### services/library_service.go
Contains business logic and data manipulation functions. Implements the LibraryManager interface in the Library struct.
