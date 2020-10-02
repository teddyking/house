/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1_test

import (
	"testing"

	"github.com/teddyking/house/api/v1alpha1"
)

var urlTests = []struct {
	name        string
	in          v1alpha1.SearchSpec
	expectedOut string
}{
	{
		"standard",
		v1alpha1.SearchSpec{
			LocationIdentifier: "REGION^475",
		},
		"https://www.rightmove.co.uk/property-for-sale/find.html?dontShow=sharedOwnership%2Cretirement&furnishTypes=&includeSSTC=false&keywords=&locationIdentifier=REGION%5E475&maxPrice=0&minBedrooms=0&mustHave=&primaryDisplayPropertyType=houses&propertyTypes=&sortType=6",
	},
	{
		"with bedrooms",
		v1alpha1.SearchSpec{
			LocationIdentifier: "REGION^475",
			MinBedrooms:        2,
		},
		"https://www.rightmove.co.uk/property-for-sale/find.html?dontShow=sharedOwnership%2Cretirement&furnishTypes=&includeSSTC=false&keywords=&locationIdentifier=REGION%5E475&maxPrice=0&minBedrooms=2&mustHave=&primaryDisplayPropertyType=houses&propertyTypes=&sortType=6",
	},
	{
		"with propertyTypes",
		v1alpha1.SearchSpec{
			LocationIdentifier: "REGION^475",
			MinBedrooms:        2,
			PropertyTypes:      "detached,semi-detached,terraced",
		},
		"https://www.rightmove.co.uk/property-for-sale/find.html?dontShow=sharedOwnership%2Cretirement&furnishTypes=&includeSSTC=false&keywords=&locationIdentifier=REGION%5E475&maxPrice=0&minBedrooms=2&mustHave=&primaryDisplayPropertyType=houses&propertyTypes=detached%2Csemi-detached%2Cterraced&sortType=6",
	},
	{
		"with everything",
		v1alpha1.SearchSpec{
			LocationIdentifier: "REGION^475",
			MinBedrooms:        2,
			MaxPrice:           500000,
			PropertyTypes:      "detached,semi-detached,terraced",
			MustHave:           "garden",
		},
		"https://www.rightmove.co.uk/property-for-sale/find.html?dontShow=sharedOwnership%2Cretirement&furnishTypes=&includeSSTC=false&keywords=&locationIdentifier=REGION%5E475&maxPrice=500000&minBedrooms=2&mustHave=garden&primaryDisplayPropertyType=houses&propertyTypes=detached%2Csemi-detached%2Cterraced&sortType=6",
	},
}

func TestURL(t *testing.T) {
	for _, tt := range urlTests {
		t.Run(tt.name, func(t *testing.T) {
			actualOut := tt.in.URL()

			if actualOut != tt.expectedOut {
				t.Errorf("got:\n %s, wanted:\n %s", actualOut, tt.expectedOut)
			}
		})
	}
}
