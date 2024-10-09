package services

import (
	"context"
	"strconv"

	"gin-books-api/cache"
	config "gin-books-api/configs"
	"gin-books-api/models"
)

// FetchAuthorsFromDB fetches authors from the database, caches them, and returns the result.
func FetchAuthorsFromDB(ctx context.Context, cacheKey string) ([]models.Author, error) {
	var authors []models.Author
	if err := config.GetDB().Preload("Book").Find(&authors).Error; err != nil {
		return nil, err
	}

	// Cache the complete list of authors
	if err := cache.SetCachedData(ctx, cacheKey, authors, cache.CacheExpiration); err != nil {
		// Proceed without caching
	}

	return authors, nil
}

// FetchAuthorFromDB fetches a single author from the database, caches it, and returns the result.
func FetchAuthorFromDB(ctx context.Context, cacheKey string, id int) (*models.Author, error) {
	var author models.Author
	if result := config.GetDB().Preload("Book").First(&author, id); result.Error != nil {
		return nil, result.Error
	}

	// Cache the fetched author
	if err := cache.SetCachedData(ctx, cacheKey, author, cache.CacheExpiration); err != nil {
		// Proceed without caching
	}

	return &author, nil
}

// CreateAuthor creates a new author and stores it in the database.
func CreateAuthor(ctx context.Context, author *models.Author) error {
	if err := config.GetDB().Create(author).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := "authors_all"
	config.RedisClient.Del(ctx, cacheKey)

	return nil
}

// UpdateAuthor updates an existing author by its ID.
func UpdateAuthor(ctx context.Context, id int, author *models.Author) error {
	author.ID = uint(id)
	if err := config.GetDB().Save(author).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := "author_" + strconv.Itoa(id)
	config.RedisClient.Del(ctx, cacheKey)
	config.RedisClient.Del(ctx, "authors_all")

	return nil
}

// DeleteAuthor deletes an author by its ID.
func DeleteAuthor(ctx context.Context, id int) error {
	if err := config.GetDB().Delete(&models.Author{}, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := "author_" + strconv.Itoa(id)
	config.RedisClient.Del(ctx, cacheKey)
	config.RedisClient.Del(ctx, "authors_all")

	return nil
}
