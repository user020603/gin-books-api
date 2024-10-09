package models

type Publisher struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`

	Books []Book `json:"books" gorm:"foreignKey:PublisherID"`
}
