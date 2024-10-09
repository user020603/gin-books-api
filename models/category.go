package models

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`

	Books []Book `json:"books" gorm:"many2many:book_categories"`
}
