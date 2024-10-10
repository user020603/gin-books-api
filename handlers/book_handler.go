package handlers

import (
	"context"
	"net/http"
	"strconv"

	"gin-books-api/cache"
	"gin-books-api/models"
	"gin-books-api/services"
	"gin-books-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetBooks retrieves all books along with their publishers, categories, authors, and reviews.
// Implements pagination and caching.
func GetBooks(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "books_all"

	// Pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	if err != nil || pageSize <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid page size")
		return
	}

	// Attempt to retrieve cached data
	var books []models.Book
	if cache.GetCachedData(ctx, cacheKey, &books) {
		// Apply pagination to cached data
		start := (page - 1) * pageSize
		end := start + pageSize
		if start > len(books) {
			start = len(books)
		}
		if end > len(books) {
			end = len(books)
		}
		paginatedBooks := books[start:end]

		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      len(books),
			"totalPages": (len(books) + pageSize - 1) / pageSize,
			"data":       paginatedBooks,
		})
		return
	}

	// If not cached, fetch from database
	books, err = services.FetchBooksFromDB(ctx, cacheKey)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve books")
		return
	}

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(books) {
		start = len(books)
	}
	if end > len(books) {
		end = len(books)
	}
	paginatedBooks := books[start:end]

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, gin.H{
		"page":       page,
		"pageSize":   pageSize,
		"total":      len(books),
		"totalPages": (len(books) + pageSize - 1) / pageSize,
		"data":       paginatedBooks,
	})
}

// GetBookByID retrieves a book by its ID along with its publisher, categories, author, and reviews.
// Implements caching and input validation.
func GetBookByID(c *gin.Context) {
	ctx := context.Background()
	idParam := c.Param("id")

	// Validate the book ID
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid book ID")
		return
	}

	cacheKey := "book_" + idParam
	var book models.Book

	// Attempt to retrieve cached data
	if cache.GetCachedData(ctx, cacheKey, &book) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, book)
		return
	}

	// If not cached, fetch from database
	bookPtr, err := services.FetchBookFromDB(ctx, cacheKey, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Book not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve the book")
		}
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, bookPtr)
}

// CreateBook creates a new book and stores it in the database.
func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.CreateBook(context.Background(), &book); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create book")
		return
	}

	utils.JSONResponse(c, http.StatusCreated, book)
}

// UpdateBook updates an existing book by its ID.
func UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.UpdateBook(context.Background(), id, &book); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update book")
		return
	}

	utils.JSONResponse(c, http.StatusOK, book)
}

// DeleteBook deletes a book by its ID.
func DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid book ID")
		return
	}

	if err := services.DeleteBook(context.Background(), id); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Book not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete book")
		}
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
