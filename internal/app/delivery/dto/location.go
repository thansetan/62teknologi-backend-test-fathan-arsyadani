package dto

type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Location struct {
	Address1       string   `json:"address1"`
	Address2       string   `json:"address2"`
	Address3       string   `json:"address3"`
	City           string   `json:"city"`
	ZipCode        string   `json:"zip_code"`
	Country        string   `json:"country"`
	State          string   `json:"state"`
	DisplayAddress []string `json:"display_address"`
}
