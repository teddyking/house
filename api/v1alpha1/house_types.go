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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HouseSpec defines the desired state of House
type HouseSpec struct {
	Price       string `json:"price"`
	OfferType   string `json:"offerType"`
	Description string `json:"description"`
	Postcode    string `json:"postcode"`
	URL         string `json:"url"`
}

// HouseStatus defines the observed state of House
type HouseStatus struct {
}

// +kubebuilder:object:root=true

// House is the Schema for the houses API
type House struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HouseSpec   `json:"spec,omitempty"`
	Status HouseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HouseList contains a list of House
type HouseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []House `json:"items"`
}

func init() {
	SchemeBuilder.Register(&House{}, &HouseList{})
}
