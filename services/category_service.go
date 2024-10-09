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
	cacheKeyCategoriesAll  = "categories_all"
	cacheKeyCategoryPrefix = "category_"
)

// FetchCategoriesFromDB fetches categories from the database, caches them, and returns the result.
func FetchCategoriesFromDB(ctx context.Context, cacheKey string) ([]models.Category, error) {
	var categories []models.Category
	if err := config.GetDB().Preload("Books").Find(&categories).Error; err != nil {
		log.Printf("Database error while fetching categories: %v", err)
		return nil, err
	}

	// Cache the complete list of authors
	if err := cache.SetCachedData(ctx, cacheKey, categories, cache.CacheExpiration); err != nil {
		log.Printf("Redis SET error for key %s: %v", cacheKey, err)
		// Proceed without caching
	}

	return categories, nil
}

func FetchCategoryFromDB(ctx context.Context, cacheKey string, id int) (*models.Category, error) {
	var category models.Category
	if result := config.GetDB().Preload("Books").First(&category, id); result.Error != nil {
		return nil, result.Error
	}

	if err := cache.SetCachedData(ctx, cacheKey, category, cache.CacheExpiration); err != nil {
		log.Printf("Failed to cache category: %v", err)
	}

	return &category, nil
}

func CreateCategory(ctx context.Context, category *models.Category) error {
	if err := config.GetDB().Create(category).Error; err != nil {
		return err
	}

	// Invalidate cache
	if err := config.RedisClient.Del(ctx, cacheKeyCategoriesAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func UpdateCategory(ctx context.Context, id int, category *models.Category) error {
	category.ID = uint(id)
	if err := config.GetDB().Save(category).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cacheKeyCategoryPrefix + strconv.Itoa(id)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
	if err := config.RedisClient.Del(ctx, cacheKeyCategoriesAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func DeleteCategory(ctx context.Context, id int) error {
	if err := config.GetDB().Delete(&models.Category{}, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cacheKeyCategoryPrefix + strconv.Itoa(id)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
	if err := config.RedisClient.Del(ctx, cacheKeyCategoriesAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}
