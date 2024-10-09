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

// GetAuthors retrieves all authors and implements caching.
func GetAuthors(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "authors_all"

	// Attempt to retrieve cached data
	var authors []models.Author
	if cache.GetCachedData(ctx, cacheKey, &authors) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, authors)
		return
	}

	// If not cached, fetch from database
	authors, err := services.FetchAuthorsFromDB(ctx, cacheKey)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve authors")
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, authors)
}

// GetAuthorByID retrieves an author by its ID and implements caching.
func GetAuthorByID(c *gin.Context) {
	ctx := context.Background()
	idParam := c.Param("id")

	// Validate the author ID
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid author ID")
		return
	}

	cacheKey := "author_" + idParam
	var author models.Author

	// Attempt to retrieve cached data
	if cache.GetCachedData(ctx, cacheKey, &author) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, author)
		return
	}

	// If not cached, fetch from database
	authorPtr, err := services.FetchAuthorFromDB(ctx, cacheKey, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Author not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve the author")
		}
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, authorPtr)
}

// CreateAuthor creates a new author and stores it in the database.
func CreateAuthor(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.CreateAuthor(context.Background(), &author); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create author")
		return
	}

	utils.JSONResponse(c, http.StatusCreated, author)
}

// UpdateAuthor updates an existing author by its ID.
func UpdateAuthor(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid author ID")
		return
	}

	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.UpdateAuthor(context.Background(), id, &author); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update author")
		return
	}

	utils.JSONResponse(c, http.StatusOK, author)
}

// DeleteAuthor deletes an author by its ID.
func DeleteAuthor(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid author ID")
		return
	}

	if err := services.DeleteAuthor(context.Background(), id); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Author not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete author")
		}
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Author deleted successfully"})
}
