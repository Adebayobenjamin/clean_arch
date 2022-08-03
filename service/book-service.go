package service

import (
	"fmt"
	"log"

	"github.com/Adebayobenjamin/clean_arch/dto"
	"github.com/Adebayobenjamin/clean_arch/entity"
	"github.com/Adebayobenjamin/clean_arch/repository"
	"github.com/mashingan/smapping"
)

// BookService is a contract that explains what this service can do
type BookService interface {
	Insert(b dto.BookCreateDTO) entity.Book
	Update(b dto.BookUpdateDTO) entity.Book
	Delete(b entity.Book)
	AllBooks() []entity.Book
	FindByID(bookID uint64) entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

// NewBookService creates an instance of the bookService struct
func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

func (service *bookService) Insert(b dto.BookCreateDTO) entity.Book {
	book := entity.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("failed map %v: ", err)
	}
	res := service.bookRepository.InsertBook(book)
	return res
}

func (service *bookService) Update(b dto.BookUpdateDTO) entity.Book {
	book := entity.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("failed map %v: ", err)
	}
	res := service.bookRepository.UpdateBook(book)
	return res
}

func (service *bookService) Delete(b entity.Book) {
	service.bookRepository.DeleteBook(b)
}

func (service *bookService) AllBooks() []entity.Book {
	return service.bookRepository.AllBooks()
}

func (service *bookService) FindByID(bookID uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookID)
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	b := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", b.UserID)
	return id == userID
}
