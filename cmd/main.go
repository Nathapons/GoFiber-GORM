package main

import (
	"backend/data"
	"backend/middleware"
	"backend/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost" // or the Docker service name if running in another container
	port     = 5432        // default PostgreSQL port
	user     = "jb"        // as defined in docker-compose.yml
	password = "12345678"  // as defined in docker-compose.yml
	dbname   = "init"      // as defined in docker-compose.yml
)

func main() {
	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Publisher{},
		&models.Author{},
		&models.AuthorBook{},
	)

	data.CreateData(db)

	// Setup Fiber
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("URL for Test main")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		userData := new(models.User)
		c.BodyParser(userData)

		token, err := models.LoginUser(db, userData)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})

		return c.Status(fiber.StatusCreated).JSON(
			fiber.Map{"message": "login success"},
		)
	})
	app.Post("/register", func(c *fiber.Ctx) error {
		userData := new(models.User)
		c.BodyParser(userData)

		if err := models.CreateUser(db, userData); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Status(fiber.StatusCreated).JSON(
			fiber.Map{"message": "Create user completed"},
		)
	})

	app.Use("/books", middleware.AuthMiddleware)

	// CRUD routes
	// app.Get("/books", func(c *fiber.Ctx) error {
	// 	return models.GetBooks(db, c)
	// })
	// app.Get("/books/:id", func(c *fiber.Ctx) error {
	// 	return models.GetBook(db, c)
	// })
	// app.Post("/books", func(c *fiber.Ctx) error {
	// 	return models.CreateBook(db, c)
	// })
	// app.Put("/books/:id", func(c *fiber.Ctx) error {
	// 	return models.UpdateBook(db, c)
	// })
	// app.Delete("/books/:id", func(c *fiber.Ctx) error {
	// 	return models.DeleteBook(db, c)
	// })

	// Start server
	log.Fatal(app.Listen(":8000"))
}
