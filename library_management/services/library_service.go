package services

import (
    "errors"
    "library_management/models"
)

type LibraryManager interface {
    AddBook(book models.Book)
    RemoveBook(bookID int)
    BorrowBook(bookID int, memberID int) error
    ReturnBook(bookID int, memberID int) error
    ListAvailableBooks() []models.Book
    ListBorrowedBooks(memberID int) []models.Book
    AddMember(member models.Member)
}

type Library struct {
    Books   map[int]models.Book
    Members map[int]models.Member
}

func NewLibrary() *Library {
    return &Library{
        Books:   make(map[int]models.Book),
        Members: make(map[int]models.Member),
    }
}

func (l *Library) AddBook(book models.Book) {
    l.Books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
    delete(l.Books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
    book, bookExists := l.Books[bookID]
    member, memberExists := l.Members[memberID]

    if !bookExists {
        return errors.New("book not found")
    }
    if !memberExists {
        return errors.New("member not found")
    }
    if book.Status == "Borrowed" {
        return errors.New("book already borrowed")
    }

    book.Status = "Borrowed"
    l.Books[bookID] = book
    member.BorrowedBooks = append(member.BorrowedBooks, book)
    l.Members[memberID] = member

    return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
    book, bookExists := l.Books[bookID]
    member, memberExists := l.Members[memberID]

    if !bookExists {
        return errors.New("book not found")
    }
    if !memberExists {
        return errors.New("member not found")
    }

    book.Status = "Available"
    l.Books[bookID] = book

    for i, b := range member.BorrowedBooks {
        if b.ID == bookID {
            member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
            break
        }
    }
    l.Members[memberID] = member

    return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
    var availableBooks []models.Book
    for _, book := range l.Books {
        if book.Status == "Available" {
            availableBooks = append(availableBooks, book)
        }
    }
    return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
    member, exists := l.Members[memberID]
    if !exists {
        return nil
    }
    return member.BorrowedBooks
}

func (l *Library) AddMember(member models.Member) {
    l.Members[member.ID] = member
}
