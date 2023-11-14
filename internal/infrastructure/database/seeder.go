package database

import (
	_ "embed"

	"gorm.io/gorm"
)

// type Category struct {
// 	Title string `json:"title"`
// 	Alias string `json:"alias"`
// }

// type Transaction struct {
// 	Type string
// }

// func getCategories(apiKey string) ([]Category, error) {
// 	var categories map[string][]Category // response is in form {"categories": [cat1....catn]}

// 	client := new(http.Client)
// 	req, err := http.NewRequest(http.MethodGet, "https://api.yelp.com/v3/categories", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
// 	req.Header.Set("Accept", "application/json")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	err = json.NewDecoder(res.Body).Decode(&categories)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return categories["categories"], nil
// }

//go:embed seeds/categories.sql
var categories string

//go:embed seeds/transactions.sql
var transactions string

func populateTable(db *gorm.DB, queries string) error {
	if err := db.Exec(queries).Error; err != nil {
		return err
	}
	return nil
}
