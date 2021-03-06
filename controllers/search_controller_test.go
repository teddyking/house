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

package controllers_test

import (
	"context"
	"errors"

	"github.com/go-logr/logr/testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/teddyking/house/api/v1alpha1"
	. "github.com/teddyking/house/controllers"
	"github.com/teddyking/house/controllers/controllersfakes"
	htypes "github.com/teddyking/house/pkg/types"
)

var _ = Describe("Search Controller", func() {
	Describe("Reconcile - Create", func() {
		var (
			fakeSearchRepo *controllersfakes.FakeSearchRepo
			fakeHouseRepo  *controllersfakes.FakeHouseRepo
			fakeScraper    *controllersfakes.FakeScraper

			namespacedName types.NamespacedName
			searchURL      string

			reconciler      *SearchReconciler
			reconcileResult ctrl.Result
			reconcileErr    error
		)

		BeforeEach(func() {
			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			fakeClient := fake.NewFakeClientWithScheme(scheme)
			fakeLogger := testing.NullLogger{}
			fakeScraper = &controllersfakes.FakeScraper{}
			fakeSearchRepo = &controllersfakes.FakeSearchRepo{}
			fakeHouseRepo = &controllersfakes.FakeHouseRepo{}

			name := "search-1"
			namespace := "namespace-1"
			namespacedName = types.NamespacedName{Name: name, Namespace: namespace}

			search := &v1alpha1.Search{
				ObjectMeta: metav1.ObjectMeta{
					Name:       name,
					Namespace:  namespace,
					Generation: 1,
				},
				Spec: v1alpha1.SearchSpec{
					LocationIdentifier: "location-id-1",
					MinBedrooms:        2,
					MaxPrice:           500000,
					PropertyTypes:      "type-1",
					MustHave:           "must-have-1",
				},
				Status: v1alpha1.SearchStatus{},
			}
			searchURL = search.Spec.URL()

			Expect(fakeClient.Create(context.TODO(), search)).To(Succeed())

			fakeSearchRepo.GetReturns(search, nil)
			fakeScraper.PropertiesReturns([]htypes.House{{Price: "price-1", Postcode: "postcode-1"}}, nil)
			fakeHouseRepo.UpsertReturns(nil)

			reconciler = &SearchReconciler{
				Client:     fakeClient,
				Log:        fakeLogger,
				Scraper:    fakeScraper,
				SearchRepo: fakeSearchRepo,
				HouseRepo:  fakeHouseRepo,
				Scheme:     scheme,
			}
		})

		JustBeforeEach(func() {
			reconcileResult, reconcileErr = reconciler.Reconcile(ctrl.Request{NamespacedName: namespacedName})
		})

		It("reconciles successfully", func() {
			Expect(reconcileErr).NotTo(HaveOccurred())
		})

		It("fetches the Search", func() {
			Expect(fakeSearchRepo.GetCallCount()).To(Equal(1))

			_, passedKey := fakeSearchRepo.GetArgsForCall(0)
			Expect(passedKey).To(Equal(namespacedName))
		})

		It("scrapes properties using the search URL", func() {
			Expect(fakeScraper.PropertiesCallCount()).To(Equal(1))

			passedURL := fakeScraper.PropertiesArgsForCall(0)
			Expect(passedURL).To(Equal(searchURL))
		})

		It("upserts House CRs from the scraped properties", func() {
			Expect(fakeHouseRepo.UpsertCallCount()).To(Equal(1))

			_, passedHouse := fakeHouseRepo.UpsertArgsForCall(0)
			Expect(passedHouse.Price).To(Equal("price-1"))
			Expect(passedHouse.Postcode).To(Equal("postcode-1"))
		})

		It("updates the Search status with the number of results", func() {
			Expect(fakeSearchRepo.UpdateStatusCallCount()).To(Equal(1))

			_, passedSearch := fakeSearchRepo.UpdateStatusArgsForCall(0)
			Expect(passedSearch.Status.NumResults).To(Equal(1))
		})

		It("updates the Search status with the observedGeneration", func() {
			Expect(fakeSearchRepo.UpdateStatusCallCount()).To(Equal(1))

			_, passedSearch := fakeSearchRepo.UpdateStatusArgsForCall(0)
			Expect(passedSearch.Status.ObservedGeneration).To(BeEquivalentTo(1))
		})

		It("reschedules another reconcile", func() {
			Expect(reconcileResult.RequeueAfter.Hours()).To(BeEquivalentTo(1))
		})

		When("there is an error fetching the Search", func() {
			BeforeEach(func() {
				fakeSearchRepo.GetReturns(nil, errors.New("error-getting-search"))
			})

			It("returns the error", func() {
				Expect(reconcileErr).To(MatchError("error-getting-search"))
			})
		})

		When("there is an error scraping the properties", func() {
			BeforeEach(func() {
				fakeScraper.PropertiesReturns([]htypes.House{}, errors.New("error-scraping-properties"))
			})

			It("returns the error", func() {
				Expect(reconcileErr).To(MatchError("error-scraping-properties"))
			})
		})

		When("there is an error creating a house", func() {
			BeforeEach(func() {
				fakeHouseRepo.UpsertReturns(errors.New("error-creating-house"))
			})

			It("returns the error", func() {
				Expect(reconcileErr).To(MatchError("error-creating-house"))
			})
		})

		When("there is an error updating the Search status", func() {
			BeforeEach(func() {
				fakeSearchRepo.UpdateStatusReturns(errors.New("error-updating-status"))
			})

			It("returns the error", func() {
				Expect(reconcileErr).To(MatchError("error-updating-status"))
			})
		})
	})
})
