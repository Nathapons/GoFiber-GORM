package data

import (
	"backend/models"

	"gorm.io/gorm"
)

func CreateData(db *gorm.DB) {
	publisher := models.Publisher{
		Details: "Publisher Details",
		Name:    "Publisher Name",
	}
	_ = models.CreatePublisher(db, &publisher)

	// Example data for a new author
	author := models.Author{
		Name: "Author Name",
	}
	_ = models.CreateAuthor(db, &author)

	// // Example data for a new book with an author
	book := models.Book{
		Name:        "Book Title",
		Author:      "Book Author",
		Description: "Book Description",
		PublisherID: publisher.ID,            // Use the ID of the publisher created above
		Authors:     []models.Author{author}, // Add the created author
	}
	_ = models.CreateBookWithAuthor(db, &book, []uint{author.ID})
}
