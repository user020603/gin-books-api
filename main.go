package main

import (
    config "gin-books-api/configs"
    "gin-books-api/handlers"
    "gin-books-api/models"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    config.InitDB()
    config.GetDB().AutoMigrate(
        &models.Author{},
        &models.BookCategory{},
        &models.Book{},
        &models.BorrowedBook{},
        &models.Category{},
        &models.Publisher{},
        &models.Review{},
        &models.User{})
    config.InitRedis()

    r := gin.Default()

    // Add CORS middleware
    r.Use(cors.Default())

    // Book routes
    r.GET("/books", handlers.GetBooks)
    r.GET("/books/:id", handlers.GetBookByID)
    r.POST("/books", handlers.CreateBook)
    r.PUT("/books/:id", handlers.UpdateBook)
    r.DELETE("/books/:id", handlers.DeleteBook)

    // Author routes
    r.GET("/authors", handlers.GetAuthors)
    r.GET("/authors/:id", handlers.GetAuthorByID)
    r.POST("/authors", handlers.CreateAuthor)
    r.PUT("/authors/:id", handlers.UpdateAuthor)
    r.DELETE("/authors/:id", handlers.DeleteAuthor)

    // Category routes
    r.GET("/categories", handlers.GetCategories)
    r.GET("/categories/:id", handlers.GetCategoryByID)
    r.POST("/categories", handlers.CreateCategory)
    r.PUT("/categories/:id", handlers.UpdateCategory)
    r.DELETE("/categories/:id", handlers.DeleteCategory)

    // Publisher routes
    r.GET("/publishers", handlers.GetPublishers)
    r.GET("/publishers/:id", handlers.GetPublisherByID)
    r.POST("/publishers", handlers.CreatePublisher)
    r.PUT("/publishers/:id", handlers.UpdatePublisher)
    r.DELETE("/publishers/:id", handlers.DeletePublisher)

    r.Run(":8080")
}