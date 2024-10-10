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
    cacheKeyUsersAll   = "users_all"
    cacheKeyUserPrefix = "user_"
)

// FetchUsersFromDB fetches users from the database, caches them, and returns the result.
func FetchUsersFromDB(ctx context.Context, cacheKey string) ([]models.User, error) {
    var users []models.User
    if err := config.GetDB().Find(&users).Error; err != nil {
        log.Printf("Database error while fetching users: %v", err)
        return nil, err
    }

    // Cache the complete list of users
    if err := cache.SetCachedData(ctx, cacheKey, users, cache.CacheExpiration); err != nil {
        log.Printf("Redis SET error for key %s: %v", cacheKey, err)
        // Proceed without caching
    }

    return users, nil
}

func FetchUserFromDB(ctx context.Context, cacheKey string, id int) (*models.User, error) {
    var user models.User
    if result := config.GetDB().First(&user, id); result.Error != nil {
        return nil, result.Error
    }

    if err := cache.SetCachedData(ctx, cacheKey, user, cache.CacheExpiration); err != nil {
        log.Printf("Failed to cache user: %v", err)
    }

    return &user, nil
}

func CreateUser(ctx context.Context, user *models.User) error {
    if err := config.GetDB().Create(user).Error; err != nil {
        return err
    }

    // Invalidate cache
    if err := config.RedisClient.Del(ctx, cacheKeyUsersAll).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }

    return nil
}

func UpdateUser(ctx context.Context, id int, user *models.User) error {
    user.ID = uint(id)
    if err := config.GetDB().Save(user).Error; err != nil {
        return err
    }

    // Invalidate cache
    cacheKey := cacheKeyUserPrefix + strconv.Itoa(id)
    if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }
    if err := config.RedisClient.Del(ctx, cacheKeyUsersAll).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }

    return nil
}

func DeleteUser(ctx context.Context, id int) error {
    if err := config.GetDB().Delete(&models.User{}, id).Error; err != nil {
        return err
    }

    // Invalidate cache
    cacheKey := cacheKeyUserPrefix + strconv.Itoa(id)
    if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }
    if err := config.RedisClient.Del(ctx, cacheKeyUsersAll).Err(); err != nil {
        log.Printf("Failed to invalidate cache: %v", err)
    }

    return nil
}