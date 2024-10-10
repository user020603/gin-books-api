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
	cacheKeyPublishersAll   = "publishers_all"
	cacheKeyPublisherPrefix = "publisher_"
)

// FetchPublishersFromDB fetches publishers from the database, caches them, and returns the result.
func FetchPublishersFromDB(ctx context.Context, cacheKey string) ([]models.Publisher, error) {
	var publishers []models.Publisher
	if err := config.GetDB().Preload("Books").Find(&publishers).Error; err != nil {
		log.Printf("Database error while fetching publishers: %v", err)
		return nil, err
	}

	// Cache the complete list of authors
	if err := cache.SetCachedData(ctx, cacheKey, publishers, cache.CacheExpiration); err != nil {
		log.Printf("Redis SET error for key %s: %v", cacheKey, err)
		// Proceed without caching
	}

	return publishers, nil
}

func FetchPublisherFromDB(ctx context.Context, cacheKey string, id int) (*models.Publisher, error) {
	var publisher models.Publisher
	if result := config.GetDB().Preload("Books").First(&publisher, id); result.Error != nil {
		return nil, result.Error
	}

	if err := cache.SetCachedData(ctx, cacheKey, publisher, cache.CacheExpiration); err != nil {
		log.Printf("Failed to cache publisher: %v", err)
	}

	return &publisher, nil
}

func CreatePublisher(ctx context.Context, publisher *models.Publisher) error {
	if err := config.GetDB().Create(publisher).Error; err != nil {
		return err
	}

	// Invalidate cache
	if err := config.RedisClient.Del(ctx, cacheKeyPublishersAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func UpdatePublisher(ctx context.Context, id int, publisher *models.Publisher) error {
	publisher.ID = uint(id)
	if err := config.GetDB().Save(publisher).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cacheKeyPublisherPrefix + strconv.Itoa(id)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
	if err := config.RedisClient.Del(ctx, cacheKeyPublishersAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func DeletePublisher(ctx context.Context, id int) error {
	if err := config.GetDB().Delete(&models.Publisher{}, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cacheKeyPublisherPrefix + strconv.Itoa(id)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
	if err := config.RedisClient.Del(ctx, cacheKeyPublishersAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}
