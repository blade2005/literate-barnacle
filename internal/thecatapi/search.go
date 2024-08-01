package thecatapi

import (
	"fmt"
	"net/http"
)

//lint:ignore U1000 stub.
type Breed struct {
	id          string
	name        string
	weight      string
	height      string
	life_span   string
	breed_group string
}

//lint:ignore U1000 stub.
type Image struct {
	id         string `json:"id"`
	url        string
	width      int
	height     int
	mime_type  string
	breeds     []Breed
	categories []string
}

type ImagesList struct {
	Count      int
	PagesCount int
	Images     []Image
}

type GetImagesSearchParams struct {
	// Size [optional] thumb , small, med or full - small is perfect for Discord
	Size *string `form:"size,omitempty" json:"size,omitempty"`

	// MimeTypes [optional] a comma separated strig of types to return e.g. jpg,png for static, or gif for gifs
	MimeTypes *string `form:"mime_types,omitempty" json:"mime_types,omitempty"`

	// Format [optional] json | src
	Format *string `form:"format,omitempty" json:"format,omitempty"`

	// HasBreeds [optional] - only return images with breed data
	HasBreeds *bool `form:"has_breeds,omitempty" json:"has_breeds,omitempty"`

	// Order [optional] default:RANDOM - RANDOM | ASC | DESC
	Order *string `form:"order,omitempty" json:"order,omitempty"`

	// Page [optional] paginate through results
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// Limit [optional] number of results to return, up to 25 with a valid API-Key
	Limit       *int    `form:"limit,omitempty" json:"limit,omitempty"`
	ContentType *string `json:"Content-Type,omitempty"`

	// XApiKey [optional] without it only the a basic set of images can be searched
	XApiKey *string `json:"x-api-key,omitempty"`
}

func (c *Client) GetImagesSearch(options *GetImagesSearchParams) (*[]Image, error) {
	search_url := fmt.Sprintf("%s/images/search", c.BaseURL)
	// TODO: setup support for params
	// My idea here is to take options and create url.Values() from it and associate it with the parsed url
	// u, err := url.Parse(search_url)
	// if err != nil {
	// 	return nil, err
	// }

	req, err := http.NewRequest("GET", search_url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	var res []Image
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
