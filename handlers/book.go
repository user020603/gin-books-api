package handlers

import (
	"encoding/json"
	"gin-books-api/config"
	"gin-books-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
    cacheKey := "books_all"

    // Check Redis cache first
    cacheBooks, err := config.RedisClient.Get(c, cacheKey).Result()
    if err == nil {
        var books []models.Book
        if err := json.Unmarshal([]byte(cacheBooks), &books); err == nil {
            c.Header("X-Data-Source", "cache")
            c.JSON(http.StatusOK, books)
            return
        }
    }

    // If not cached
    var books []models.Book
    if result := config.GetDB().Find(&books); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    booksJSON, _ := json.Marshal(books)
    config.RedisClient.Set(c, cacheKey, booksJSON, 10*time.Minute)

    c.Header("X-Data-Source", "database")
    c.JSON(http.StatusOK, books)
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "book_" + id

	// Check Redis cache first
	cacheBook, err := config.RedisClient.Get(c, cacheKey).Result()
	if err == nil {
		var book models.Book
		if err := json.Unmarshal([]byte(cacheBook), &book); err == nil {
			c.Header("X-Data-Source", "cache")
			c.JSON(http.StatusOK, book)
			return
		}
	}

	// If not cached
	var book models.Book
	if result := config.GetDB().First(&book, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	bookJSON, _ := json.Marshal(book)
	config.RedisClient.Set(c, cacheKey, bookJSON, 10*time.Minute)

	c.Header("X-Data-Source", "database")
	c.JSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate book
	var existingBook models.Book
	if result := config.GetDB().Where("title = ?", book.Title).First(&existingBook); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Book with the same title already exists"})
		return
	}

	if result := config.GetDB().Create(&book); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// add new data, cache is changed, so empty cache
	config.RedisClient.Del(c, "books_all")

	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if result := config.GetDB().First(&book, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.GetDB().Save(&book)

	// Invalidate cache
	config.RedisClient.Del(c, "books_all", "book_"+id)

	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if result := config.GetDB().First(&book, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	config.GetDB().Delete(&book)

	config.RedisClient.Del(c, "books_all", "book_"+id)

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
