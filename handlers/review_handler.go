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

func GetReviews(c *gin.Context) {
    ctx := context.Background()
    cacheKey := "reviews_all"

    var reviews []models.Review
    if cache.GetCachedData(ctx, cacheKey, &reviews) {
        c.Header("X-Data-Source", "cache")
        utils.JSONResponse(c, http.StatusOK, reviews)
        return
    }

    reviews, err := services.FetchReviewsFromDB(ctx, cacheKey)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reviews")
        return
    }

    c.Header("X-Data-Source", "database")
    utils.JSONResponse(c, http.StatusOK, reviews)
}

func GetReviewByID(c *gin.Context) {
    ctx := context.Background()
    idParam := c.Param("id")

    id, err := strconv.Atoi(idParam)
    if err != nil || id <= 0 {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
        return
    }

    cacheKey := "review_" + idParam
    var review models.Review

    if cache.GetCachedData(ctx, cacheKey, &review) {
        c.Header("X-Data-Source", "cache")
        utils.JSONResponse(c, http.StatusOK, review)
        return
    }

    reviewPtr, err := services.FetchReviewFromDB(ctx, cacheKey, id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.ErrorResponse(c, http.StatusNotFound, "Review not found")
        } else {
            utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve the review")
        }
        return
    }

    c.Header("X-Data-Source", "database")
    utils.JSONResponse(c, http.StatusOK, reviewPtr)
}

func CreateReview(c *gin.Context) {
    var review models.Review
    if err := c.ShouldBindJSON(&review); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
        return
    }

    if err := services.CreateReview(context.Background(), &review); err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create review")
        return
    }

    utils.JSONResponse(c, http.StatusCreated, review)
}

func UpdateReview(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil || id <= 0 {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
        return
    }

    var review models.Review
    if err := c.ShouldBindJSON(&review); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
        return
    }

    if err := services.UpdateReview(context.Background(), id, &review); err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update review")
        return
    }

    utils.JSONResponse(c, http.StatusOK, review)
}

func DeleteReview(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil || id <= 0 {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID")
        return
    }

    if err := services.DeleteReview(context.Background(), id); err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.ErrorResponse(c, http.StatusNotFound, "Review not found")
        } else {
            utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete review")
        }
        return
    }

    utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Review deleted successfully"})
}