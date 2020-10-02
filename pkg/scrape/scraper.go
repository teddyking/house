package scrape

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/teddyking/house/pkg/types"
)

type RightMove struct {
	Collector *colly.Collector
}

type Result struct {
	Properties []Property `json:"properties"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type Property struct {
	ID             int      `json:"id"`
	Type           string   `json:"propertySubType"`
	Description    string   `json:"propertyTypeFullDescription"`
	Bedrooms       int      `json:"bedrooms"`
	Location       Location `json:"location"`
	DisplayAddress string   `json:"displayAddress"`
	Price          Price    `json:"price"`
	URL            string   `json:"propertyUrl"`
}

type Pagination struct {
	Total int    `json:"total"`
	First string `json:"first"`
	Last  string `json:"last"`
	Next  string `json:"next"`
	Page  string `json:"page"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Price struct {
	Amount       int    `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

// Properties is a fairly simple and unintelligent func ... It's a bit of a
// hack and is currently completely untested. V. likely to break.
func (r RightMove) Properties(url string) ([]types.House, error) {
	searchResultPagination, err := r.getSearchResultPagination(url)
	if err != nil {
		return []types.House{}, err
	}

	searchResultURLs, err := getSearchResultURLs(url, searchResultPagination)
	if err != nil {
		return []types.House{}, err
	}

	searchResultJSONs, err := r.getSearchResultJSONs(searchResultURLs)
	if err != nil {
		return []types.House{}, err
	}

	properties := []types.House{}
	for _, searchResultJSON := range searchResultJSONs {
		result := &Result{}
		if err := json.NewDecoder(strings.NewReader(searchResultJSON)).Decode(result); err != nil {
			return []types.House{}, err
		}

		for _, property := range result.Properties {
			properties = append(properties, types.House{
				ID:          fmt.Sprint(property.ID),
				Price:       fmt.Sprint(property.Price.Amount),
				Description: property.Description,
				Postcode:    postcodeFromAddress(property.DisplayAddress),
				URL:         property.URL,
			})
		}
	}

	return properties, nil
}

func (r RightMove) getSearchResultPagination(url string) (Pagination, error) {
	paginationCollector := r.Collector.Clone()

	paginationJSONRaw := ""
	paginationCollector.OnHTML("script", func(element *colly.HTMLElement) {
		text := element.Text

		if strings.HasPrefix(text, "window.jsonModel =") {
			paginationJSONRaw = strings.TrimPrefix(text, "window.jsonModel = ")
		}
	})

	if err := paginationCollector.Visit(url); err != nil {
		return Pagination{}, err
	}

	paginationResult := &Result{}
	if err := json.NewDecoder(strings.NewReader(paginationJSONRaw)).Decode(paginationResult); err != nil {
		return Pagination{}, err
	}

	return paginationResult.Pagination, nil
}

func getSearchResultURLs(url string, pagination Pagination) ([]string, error) {
	searchResultURLs := []string{}
	if pagination.Total <= 1 {
		searchResultURLs = []string{url}
	}

	if pagination.Total > 1 {
		increment, err := strconv.Atoi(pagination.Next)
		if err != nil {
			return []string{}, err
		}
		last, _ := strconv.Atoi(pagination.Last)
		if err != nil {
			return []string{}, err
		}
		numPages := (last / increment) + 1

		for i := 0; i < numPages; i++ {
			index := i * increment
			urlWithIndex := fmt.Sprintf("%s&index=%d", url, index)
			searchResultURLs = append(searchResultURLs, urlWithIndex)
		}
	}

	return searchResultURLs, nil
}

func (r RightMove) getSearchResultJSONs(searchResultURLs []string) ([]string, error) {
	propertyCollector := r.Collector.Clone()

	searchResultJSONs := map[string]string{}
	propertyCollector.OnHTML("script", func(element *colly.HTMLElement) {
		text := element.Text

		if strings.HasPrefix(text, "window.jsonModel =") {
			searchResultJSONRaw := strings.TrimPrefix(text, "window.jsonModel = ")
			searchResultJSONs[element.Request.URL.String()] = searchResultJSONRaw
		}
	})

	for _, propertyURL := range searchResultURLs {
		if err := propertyCollector.Visit(propertyURL); err != nil {
			return []string{}, err
		}
	}

	resultJSONs := []string{}
	for _, json := range searchResultJSONs {
		resultJSONs = append(resultJSONs, json)
	}

	return resultJSONs, nil
}

func postcodeFromAddress(address string) string {
	return strings.TrimSpace(address[strings.LastIndex(address, ",")+1:])
}
