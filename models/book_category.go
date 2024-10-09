package models

type BookCategory struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	BookID uint `json:"book_id" gorm:"not null"`
	CategoryID uint `json:"category_id" gorm:"not null"`

	Book Book `json:"book" gorm:"foreignKey:BookID"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
}
