package handlers

import (
	"context"
	"gin-books-api/cache"
	"gin-books-api/models"
	"gin-books-api/services"
	"gin-books-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCategories retrieves all categories and implements caching.
func GetCategories(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "categories_all"

	// Attempt to retrieve cached data
	var categories []models.Category
	if cache.GetCachedData(ctx, cacheKey, &categories) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, categories)
		return
	}

	// If not cached, fetch from database
	categories, err := services.FetchCategoriesFromDB(ctx, cacheKey)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, categories)
}

// GetCategoryByID retrieves an category by its ID and implements caching.
func GetCategoryByID(c *gin.Context) {
	ctx := context.Background()
	idParam := c.Param("id")

	// Validate the categories ID
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	cacheKey := "category_" + idParam
	var category models.Category

	// Attempt to retrieve cached data
	if cache.GetCachedData(ctx, cacheKey, &category) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, category)
		return
	}

	// If not cached, fetch from database
	categoryPtr, err := services.FetchCategoryFromDB(ctx, cacheKey, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve the category")
		}
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, categoryPtr)
}

// CreateCategory creates a new category and stores it in the database.
func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.CreateCategory(context.Background(), &category); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category")
		return
	}

	utils.JSONResponse(c, http.StatusCreated, category)
}

// UpdateCategory updates an existing category by its ID.
func UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.UpdateCategory(context.Background(), id, &category); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category")
		return
	}

	utils.JSONResponse(c, http.StatusOK, category)
}

// DeleteCategory deletes an category by its ID.
func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := services.DeleteCategory(context.Background(), id); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category")
		}
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
