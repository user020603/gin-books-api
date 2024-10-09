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
	cacheKeyBooksAll   = "books_all"
	cacheKeyBookPrefix = "book_"
)

// FetchBooksFromDB fetches books from the database, caches them, and returns the result.
func FetchBooksFromDB(ctx context.Context, cacheKey string) ([]models.Book, error) {
	var books []models.Book
	query := config.GetDB().
		Preload("Author"). // Preload books for each author
		Preload("Publisher").
		Preload("Categories").
		Preload("Reviews")

	// Fetch all books for caching
	if err := query.Find(&books).Error; err != nil {
		log.Printf("Database error while fetching books: %v", err)
		return nil, err
	}

	// Cache the complete list of books
	if err := cache.SetCachedData(ctx, cacheKey, books, cache.CacheExpiration); err != nil {
		log.Printf("Redis SET error for key %s: %v", cacheKey, err)
		// Proceed without caching
	}

	return books, nil
}

// FetchBookFromDB fetches a single book from the database, caches it, and returns the result.
func FetchBookFromDB(ctx context.Context, cacheKey string, id int) (*models.Book, error) {
	var book models.Book
	if result := config.GetDB().Preload("Author").Preload("Publisher").Preload("Categories").Preload("Reviews").First(&book, id); result.Error != nil {
		return nil, result.Error
	}

	// Cache the fetched book
	if err := cache.SetCachedData(ctx, cacheKey, book, cache.CacheExpiration); err != nil {
		log.Printf("Redis SET error for key %s: %v", cacheKey, err)
		// Proceed without caching
	}

	return &book, nil
}

// CreateBook creates a new book and stores it in the database.
func CreateBook(ctx context.Context, book *models.Book) error {
	if err := config.GetDB().Create(book).Error; err != nil {
		return err
	}

	// Invalidate cache
	if err := config.RedisClient.Del(ctx, cacheKeyBooksAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

// UpdateBook updates an existing book by its ID.
func UpdateBook(ctx context.Context, id int, book *models.Book) error {
	book.ID = uint(id)
	if err := config.GetDB().Save(book).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cacheKeyBookPrefix + strconv.Itoa(id)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
	if err := config.RedisClient.Del(ctx, cacheKeyBooksAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

// DeleteBook deletes a book by its ID.
func DeleteBook(ctx context.Context, id int) error {
	if err := config.GetDB().Delete(&models.Book{}, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cacheKeyBookPrefix + strconv.Itoa(id)
	if err := config.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
	if err := config.RedisClient.Del(ctx, cacheKeyBooksAll).Err(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}
