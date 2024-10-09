package services

import (
    "context"
    "log"
    "strconv"

    "gin-books-api/cache"
    config "gin-books-api/configs"
    "gin-books-api/models"
)

const (
    cacheKeyAuthorsAll     = "authors_all"
    cacheKeyAuthorPrefix   = "author_"
)

// FetchAuthorsFromDB fetches authors from the database, caches them, and returns the result.
func FetchAuthorsFromDB(ctx context.Context, cacheKey string) ([]models.Author, error) {
    var authors []models.Author
    if err := config.GetDB().Preload("Book").Find(&authors).Error; err != nil {
        log.Printf("Database error while fetching authors: %v", err)
        return nil, err
    }

    // Cache the complete list of authors
    if err := cache.SetCachedData(ctx, cacheKey, authors, cache.CacheExpiration); err != nil {
        log.Printf("Redis SET error for key %s: %v", cacheKey, err)
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
        log.Printf("Redis SET error for key %s: %v", cacheKey, err)
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
    if err := config.RedisClient.Del(ctx, cacheKeyAuthorsAll).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }

    return nil
}

// UpdateAuthor updates an existing author by its ID.
func UpdateAuthor(ctx context.Context, id int, author *models.Author) error {
    author.ID = uint(id)
    if err := config.GetDB().Save(author).Error; err != nil {
        return err
    }

    // Invalidate cache
    cacheKey := cacheKeyAuthorPrefix + strconv.Itoa(id)
    if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }
    if err := config.RedisClient.Del(ctx, cacheKeyAuthorsAll).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }

    return nil
}

// DeleteAuthor deletes an author by its ID.
func DeleteAuthor(ctx context.Context, id int) error {
    if err := config.GetDB().Delete(&models.Author{}, id).Error; err != nil {
        return err
    }

    // Invalidate cache
    cacheKey := cacheKeyAuthorPrefix + strconv.Itoa(id)
    if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }
    if err := config.RedisClient.Del(ctx, cacheKeyAuthorsAll).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }

    return nil
}