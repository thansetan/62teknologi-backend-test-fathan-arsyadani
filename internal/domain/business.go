package domain

type Business struct {
	ID              string `gorm:"primaryKey"`
	Name            string
	Alias           string `gorm:"index"`
	ImageURL        string
	Price           string
	Phone           string
	OpenAt, CloseAt string `gorm:"index"`

	Categories   []Category    `gorm:"many2many:businesses_categories;constraint:onDelete:CASCADE"`
	Transactions []Transaction `gorm:"many2many:businesses_transactions;constraint:onDelete:CASCADE"`
	Location     Location      `gorm:"foreignKey:BusinessID;references:ID;constraint:onDelete:CASCADE"`
}

type BusinessQuery struct {
	Location        string
	Categories      []string
	Limit           int
	Offset          int
	TransactionType []string
	OpenNow         string
	OpenAt          string
}
