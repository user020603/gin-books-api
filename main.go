package main

import (
	"gin-books-api/config"
	"gin-books-api/handlers"
	"gin-books-api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	config.GetDB().AutoMigrate(&models.Book{})

	r := gin.Default()

	r.GET("/books", handlers.GetBooks)
	r.GET("/books/:id", handlers.GetBookByID)
	r.POST("/books", handlers.CreateBook)
	r.PUT("/books/:id", handlers.UpdateBook)
	r.DELETE("/books/:id", handlers.DeleteBook)

	r.Run(":8080")
}
