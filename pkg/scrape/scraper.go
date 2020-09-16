package scrape

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/teddyking/house/pkg/types"
)

type RightMove struct {
	Collector *colly.Collector
}

type Result struct {
	Properties []Property `json:"properties"`
}

type Property struct {
	ID          int      `json:"id"`
	Type        string   `json:"propertySubType"`
	Description string   `json:"propertyTypeFullDescription"`
	Bedrooms    int      `json:"bedrooms"`
	Location    Location `json:"location"`
	Price       Price    `json:"price"`
	URL         string   `json:"propertyUrl"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Price struct {
	Amount       int    `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

func (r RightMove) Properties(url string) ([]types.House, error) {
	propertyJSONRaw := ""

	r.Collector.OnHTML("script", func(element *colly.HTMLElement) {
		text := element.Text

		if strings.HasPrefix(text, "window.jsonModel =") {
			propertyJSONRaw = strings.TrimPrefix(text, "window.jsonModel = ")
		}
	})

	if err := r.Collector.Visit(url); err != nil {
		return []types.House{}, err
	}

	result := &Result{}
	if err := json.NewDecoder(strings.NewReader(propertyJSONRaw)).Decode(result); err != nil {
		return []types.House{}, err
	}

	houses := []types.House{}
	for _, property := range result.Properties {
		houses = append(houses, types.House{
			ID:          fmt.Sprint(property.ID),
			Price:       fmt.Sprint(property.Price.Amount),
			Description: property.Description,
			URL:         property.URL,
		})
	}

	return houses, nil
}
