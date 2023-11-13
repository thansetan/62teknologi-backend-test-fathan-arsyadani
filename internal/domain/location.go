package domain

type Location struct {
	ID                  uint `gorm:"primaryKey"`
	Address1            string
	Address2            string
	Address3            string
	City                string
	ZipCode             int
	Country             string
	State               string
	Latitude, Longitude float64

	BusinessID string
}
