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

func GetUsers(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "users_all"

	var users []models.User
	if cache.GetCachedData(ctx, cacheKey, &users) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, users)
		return
	}

	users, err := services.FetchUsersFromDB(ctx, cacheKey)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	ctx := context.Background()
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	cacheKey := "user_" + idParam
	var user models.User

	if cache.GetCachedData(ctx, cacheKey, &user) {
		c.Header("X-Data-Source", "cache")
		utils.JSONResponse(c, http.StatusOK, user)
		return
	}

	userPtr, err := services.FetchUserFromDB(ctx, cacheKey, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve the user")
		}
		return
	}

	c.Header("X-Data-Source", "database")
	utils.JSONResponse(c, http.StatusOK, userPtr)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.CreateUser(context.Background(), &user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.JSONResponse(c, http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := services.UpdateUser(context.Background(), id, &user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.JSONResponse(c, http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := services.DeleteUser(context.Background(), id); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}
