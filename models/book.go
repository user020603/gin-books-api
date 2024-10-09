package models

type Book struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	PublishedYear int        `json:"published_year"`
	AuthorID      *uint      `json:"author_id"`    // Use pointer to allow NULL values
	PublisherID   *uint      `json:"publisher_id"` // Use pointer to allow NULL values
	Availability   bool      `json:"availability" gorm:"default:true"` // Indicates if the book is available for borrowing

	Author     Author     `json:"-" gorm:"foreignKey:AuthorID"`   // Relation to Author
	Publisher  Publisher  `json:"-" gorm:"foreignKey:PublisherID"` // Relation to Publisher
	Reviews    []Review   `json:"-" gorm:"foreignKey:BookID"`      // Relation to Reviews
	Categories []Category `json:"-" gorm:"many2many:book_categories;"` // Many-to-many relationship with Categories
}