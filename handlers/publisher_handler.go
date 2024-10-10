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

func GetPublishers(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "publishers_all"

	var publishers []models.Publisher
	if cache.GetCachedData(ctx, cacheKey, &publishers) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, publishers)
		return
	}

	publishers, err := services.FetchPublishersFromDB(ctx, cacheKey)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve publishers")
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, publishers)
}

func GetPublisherByID(c *gin.Context) {
	ctx := context.Background()
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid publisher ID")
		return
	}

	cacheKey := "publisher_" + idParam
	var publisher models.Publisher

	if cache.GetCachedData(ctx, cacheKey, &publisher) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, publisher)
		return
	}

	publisherPtr, err := services.FetchPublisherFromDB(ctx, cacheKey, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Publisher not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve the publisher")
		}
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, publisherPtr)
}

func CreatePublisher(c *gin.Context) {
	var publisher models.Publisher
	if err := c.ShouldBindJSON(&publisher); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.CreatePublisher(context.Background(), &publisher); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create publisher")
		return
	}

	utils.JSONResponse(c, http.StatusCreated, publisher)
}

func UpdatePublisher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid publisher ID")
		return
	}

	var publisher models.Publisher
	if err := c.ShouldBindJSON(&publisher); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.UpdatePublisher(context.Background(), id, &publisher); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update publisher")
		return
	}

	utils.JSONResponse(c, http.StatusOK, publisher)
}

func DeletePublisher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid publisher ID")
		return
	}

	if err := services.DeletePublisher(context.Background(), id); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Publisher not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete publisher")
		}
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Publisher deleted successfully"})
}
