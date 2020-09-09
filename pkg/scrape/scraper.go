package scrape

import "github.com/teddyking/house/pkg/types"

type RightMove struct{}

func (r RightMove) Properties(url string) ([]types.House, error) {
	return []types.House{
		{
			ID:          "112233",
			Price:       "350,000",
			OfferType:   "Offers Over",
			Postcode:    "EH1",
			Description: "A great house",
			URL:         "url1",
		},
		{
			ID:          "445566",
			Price:       "500,000",
			OfferType:   "Offers Over",
			Postcode:    "EH2",
			Description: "An even better house",
			URL:         "url2",
		},
	}, nil
}
