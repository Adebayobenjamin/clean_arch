package repository

import (
	"github.com/Adebayobenjamin/clean_arch/entity"
	"gorm.io/gorm"
)

// BookRepository is a contract that specifies what this repository can do
type BookRepository interface {
	InsertBook(b entity.Book) entity.Book
	UpdateBook(b entity.Book) entity.Book
	DeleteBook(b entity.Book)
	AllBooks() []entity.Book
	FindBookByID(bookID uint64) entity.Book
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookConnection{
		connection: db,
	}
}

func (db *bookConnection) InsertBook(b entity.Book) entity.Book {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *bookConnection) UpdateBook(b entity.Book) entity.Book {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *bookConnection) DeleteBook(b entity.Book) {
	db.connection.Delete(&b)
}

func (db *bookConnection) AllBooks() []entity.Book {
	var books []entity.Book
	db.connection.Preload("User").Find(&books)
	return books
}

func (db *bookConnection) FindBookByID(bookID uint64) entity.Book {
	var book entity.Book
	db.connection.Preload("User").Find(&book, bookID)
	return book
}
