package domain

type Category struct {
	ID    uint `gorm:"primaryKey"`
	Title string
	Alias string `gorm:"index"`
}
