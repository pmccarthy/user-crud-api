package main

import (
	"log"
	"net/http"
	"os"
	"user-crud-api/models"
	"user-crud-api/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Database struct {
	DB *gorm.DB
}

func (d *Database) CreateUser(context *fiber.Ctx) error {
	user := User{}
	err := context.BodyParser(&user)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Request failed",
		})
		return err
	}

	err = d.DB.Create(&user).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not create user",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{})

	log.Printf("Successfully created user: %s email: %s", user.Username, user.Email)

	return nil
}

func (d *Database) DeleteUser(context *fiber.Ctx) error {
	userModel := models.Users{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id parameter cannot be empty",
		})
		return nil
	}

	err := d.DB.Delete(userModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not delete user",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{})

	log.Printf("Deleted user with id %s successfully", id)

	return nil
}

func (d *Database) GetUsers(context *fiber.Ctx) error {
	userModels := &[]models.Users{}

	err := d.DB.Find(userModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"data": userModels,
	})

	log.Println("Get users request successful")

	return nil
}

func (d *Database) GetUserbyID(context *fiber.Ctx) error {
	userModel := &models.Users{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := d.DB.Where("id = ?", id).First(userModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not find user with specified ID",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"data": userModel,
	})

	log.Printf("Retrieved user by ID: %s", id)

	return nil
}

func (r *Database) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_user", r.CreateUser)
	api.Get("/get_user/:id", r.GetUserbyID)
	api.Get("/get_users", r.GetUsers)
	api.Delete("/delete_user/:id", r.DeleteUser)
	//TODO: implement update endpoint
	// api.Post("/update_user", r.UpdateUser)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}
	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("Could not migrate db")
	}

	d := Database{
		DB: db,
	}
	app := fiber.New()
	d.SetupRoutes(app)
	app.Listen(":8080")
	app.Use(logger.New())
}
