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

package v1alpha1

import (
	"fmt"
	"net/url"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Search is the Schema for the searches API
type Search struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SearchSpec   `json:"spec,omitempty"`
	Status SearchStatus `json:"status,omitempty"`
}

// SearchSpec defines the desired state of Search
type SearchSpec struct {
	LocationIdentifier string `json:"locationIdentifier"`
	MinBedrooms        int    `json:"minBedrooms,omitempty"`
	MaxPrice           int    `json:"maxPrice,omitempty"`
	PropertyTypes      string `json:"propertyTypes,omitempty"`
	MustHave           string `json:"mustHave,omitempty"`
}

func (s SearchSpec) URL() string {
	v := url.Values{}

	v.Set("locationIdentifier", s.LocationIdentifier)
	v.Add("minBedrooms", strconv.Itoa(s.MinBedrooms))
	v.Add("maxPrice", strconv.Itoa(s.MaxPrice))
	v.Add("propertyTypes", s.PropertyTypes)
	v.Add("mustHave", s.MustHave)

	v.Add("sortType", "6")
	v.Add("primaryDisplayPropertyType", "houses")
	v.Add("includeSSTC", "false")
	v.Add("dontShow", "sharedOwnership,retirement")
	v.Add("furnishTypes", "")
	v.Add("keywords", "")

	return fmt.Sprintf("https://www.rightmove.co.uk/property-for-sale/find.html?%s", v.Encode())
}

// SearchStatus defines the observed state of Search
type SearchStatus struct {
	ObservedGeneration int64 `json:"observedGeneration"`
	NumResults         int   `json:"numResults"`
}

// +kubebuilder:object:root=true
// SearchList contains a list of Search
type SearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Search `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Search{}, &SearchList{})
}
