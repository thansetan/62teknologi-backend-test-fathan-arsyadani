package domain

type Transaction struct {
	ID   uint   `gorm:"primaryKey"`
	Type string `gorm:"index"`
}
