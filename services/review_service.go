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
	cacheKeyReviewsAll   = "reviews_all"
	cacheKeyReviewPrefix = "review_"
)

// FetchReviewsFromDB fetches reviews from the database, caches them, and returns the result.
func FetchReviewsFromDB(ctx context.Context, cacheKey string) ([]models.Review, error) {
	var reviews []models.Review
	if err := config.GetDB().Preload("Book").Preload("User").Find(&reviews).Error; err != nil {
		log.Printf("Database error while fetching reviews: %v", err)
		return nil, err
	}

	// Cache the complete list of reviews
	if err := cache.SetCachedData(ctx, cacheKey, reviews, cache.CacheExpiration); err != nil {
		log.Printf("Redis SET error for key %s: %v", cacheKey, err)
		// Proceed without caching
	}

	return reviews, nil
}

func FetchReviewFromDB(ctx context.Context, cacheKey string, id int) (*models.Review, error) {
	var review models.Review
	if result := config.GetDB().Preload("Book").Preload("User").First(&review, id); result.Error != nil {
		return nil, result.Error
	}

	if err := cache.SetCachedData(ctx, cacheKey, review, cache.CacheExpiration); err != nil {
		log.Printf("Failed to cache review: %v", err)
	}

	return &review, nil
}

func CreateReview(ctx context.Context, review *models.Review) error {
	if err := config.GetDB().Create(review).Error; err != nil {
		return err
	}

	// Invalidate cache
	if err := config.RedisClient.Del(ctx, cacheKeyReviewsAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func UpdateReview(ctx context.Context, id int, review *models.Review) error {
	review.ID = uint(id)
	if err := config.GetDB().Save(review).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cacheKeyReviewPrefix + strconv.Itoa(id)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
	if err := config.RedisClient.Del(ctx, cacheKeyReviewsAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func DeleteReview(ctx context.Context, id int) error {
	if err := config.GetDB().Delete(&models.Review{}, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cacheKeyReviewPrefix + strconv.Itoa(id)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
	if err := config.RedisClient.Del(ctx, cacheKeyReviewsAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}
