package dto

import (
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type BusinessRequest struct {
	Name         string  `json:"name" example:"John Doe John Store"`
	ImageURL     string  `json:"image_url" example:"https://john.doe/image.png"`
	Categories   *string `json:"categories" example:"venues,adult,antiques"` // comma separated ex: (japanese,venues,etc). to make it empty, just do ""
	Latitude     float64 `json:"latitude" example:"37.234332396"`
	Longitude    float64 `json:"longitude" example:"-115.80666344"`
	Transactions *string `json:"transactions" example:"pickup,delivery"` // comma separated ex: (pickup,delivery,etc). to make it empty, just do ""
	Price        string  `json:"price" example:"$$$$"`
	Address1     string  `json:"address1" example:"53rd John Doe Street"`
	Address2     string  `json:"address2" example:"Johnnyslvania"`
	Address3     string  `json:"address3" example:""`
	City         string  `json:"city" example:"Lincoln"`
	ZipCode      int     `json:"zip_code" example:"68501"`
	Country      string  `json:"country" example:"US"`
	State        string  `json:"state" example:"Nebraska"`
	Phone        string  `json:"phone" example:"+1938485823"`
	OpenAt       string  `json:"open_at" example:"1100"`
	CloseAt      string  `json:"close_at" example:"2100"`
}

func (b BusinessRequest) ValidateCreate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name, validation.Required),
		validation.Field(&b.ImageURL, validation.Required, is.URL),
		validation.Field(&b.Latitude, validation.By(isLatitude)),
		validation.Field(&b.Longitude, validation.By(isLongitude)),
		validation.Field(&b.Price, validation.Required),
		validation.Field(&b.Address1, validation.Required),
		validation.Field(&b.City, validation.Required),
		validation.Field(&b.ZipCode, validation.Required),
		validation.Field(&b.Country, validation.Required),
		validation.Field(&b.State, validation.Required),
		validation.Field(&b.Phone, validation.Required, is.E164),
		validation.Field(&b.OpenAt, validation.By(isHour)),
		validation.Field(&b.CloseAt, validation.By(isHour)),
	)
}

func (b BusinessRequest) ValidateUpdate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.ImageURL, is.URL),
		validation.Field(&b.Latitude, validation.By(isLatitude)),
		validation.Field(&b.Longitude, validation.By(isLongitude)),
		validation.Field(&b.Phone, is.E164),
		validation.Field(&b.OpenAt, validation.When(b.OpenAt != "", validation.By(isHour))),
		validation.Field(&b.CloseAt, validation.When(b.CloseAt != "", validation.By(isHour))),
	)
}

type BusinessQueryParams struct {
	Location        string `query:"location"`   // city/country/state/address
	Categories      string `query:"categories"` // comma separated(cat1,cat2,...,catn)
	Page            int    `query:"page"`
	PerPage         int    `query:"per_page"`
	TransactionType string `query:"transactions"` // comma separated(cat1,cat2,...,catn)
	OpenNow         bool   `query:"open_now"`     // when set to false or not set, will return all
	OpenAt          string `query:"open_at"`
}

func (p BusinessQueryParams) CategoriesList() []string {
	return strings.Split(p.Categories, ",")
}

func (p BusinessQueryParams) TransactionsList() []string {
	return strings.Split(p.TransactionType, ",")
}

func (b BusinessQueryParams) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Location, validation.Required.Error("please specity a location")),
		validation.Field(&b.OpenNow, validation.When(b.OpenAt != "", validation.Empty.Error("can't use both open_now and open_at"))),
		validation.Field(&b.OpenAt, validation.When(b.OpenAt != "", validation.By(isHour)), validation.When(b.OpenNow, validation.Empty.Error("can't use both open_now and open_at"))),
	)
}

type BusinessResponse struct {
	ID           string      `json:"id"`
	Alias        string      `json:"alias"`
	Name         string      `json:"name"`
	ImageURL     string      `json:"imageURL"`
	IsOpen       *bool       `json:"is_open,omitempty"`
	URL          string      `json:"url"`
	Categories   []Category  `json:"categories"`
	Coordinates  Coordinates `json:"coordinates"`
	Transactions []string    `json:"transactions"`
	Price        string      `json:"price"`
	Location     Location    `json:"location"`
	Phone        string      `json:"phone"`
	OpenAt       string      `json:"open_at,omitempty"`
	CloseAt      string      `json:"close_at,omitempty"`
}

type BusinessesResponse struct {
	Businesses []BusinessResponse
	Metadata   Metadata
}

type Metadata struct {
	Total       int64     `json:"total"`
	PerPage     int       `json:"per_page"`
	Page        int       `json:"page"`
	TotalPages  int       `json:"total_pages"`
	CurrentTime time.Time `json:"current_time"`
}

func isHour(value any) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("must be string")
	}
	if len(s) != 4 || s < "0000" || s > "2359" {
		return errors.New(`time must be between "0000" and "2359"`)
	}
	return nil
}

func isLatitude(value any) error {
	l, ok := value.(float64)
	if !ok {
		return errors.New("must be numeric")
	}

	if l < -90 || l > 90 {
		return errors.New("must be between -90 and 90")
	}

	return nil
}

func isLongitude(value any) error {
	l, ok := value.(float64)
	if !ok {
		return errors.New("must be numeric")
	}

	if l < -180 || l > 180 {
		return errors.New("must be between -180 and 180")
	}

	return nil
}
