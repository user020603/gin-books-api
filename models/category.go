package models

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`

	Books []Book `json:"books" gorm:"foreignKey:CategoryID"`
}
