package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

// getBooks retrieves all books
func GetBooks(db *gorm.DB, c *fiber.Ctx) error {
	var books []Book
	db.Find(&books)
	return c.JSON(books)
}

// getBook retrieves a book by id
func GetBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	db.First(&book, id)
	return c.JSON(book)
}

// createBook creates a new book
func CreateBook(db *gorm.DB, c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return err
	}
	db.Create(&book)
	return c.JSON(book)
}

// updateBook updates a book by id
func UpdateBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	book := new(Book)
	db.First(&book, id)
	if err := c.BodyParser(book); err != nil {
		return err
	}
	db.Save(&book)
	return c.JSON(book)
}

// deleteBook deletes a book by id
func DeleteBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&Book{}, id)
	return c.SendString("Book successfully deleted")
}
