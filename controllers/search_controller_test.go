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

			reconciler   *SearchReconciler
			reconcileErr error
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
					URL: "url-1",
				},
				Status: v1alpha1.SearchStatus{},
			}

			Expect(fakeClient.Create(context.TODO(), search)).To(Succeed())

			fakeSearchRepo.GetReturns(search, nil)
			fakeScraper.PropertiesReturns([]htypes.House{{Price: "price-1", Postcode: "postcode-1"}}, nil)
			fakeHouseRepo.CreateReturns(nil)

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
			_, reconcileErr = reconciler.Reconcile(ctrl.Request{NamespacedName: namespacedName})
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
			Expect(passedURL).To(Equal("url-1"))
		})

		It("creates House CRs from the scraped properties", func() {
			Expect(fakeHouseRepo.CreateCallCount()).To(Equal(1))

			_, passedHouse := fakeHouseRepo.CreateArgsForCall(0)
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
				fakeHouseRepo.CreateReturns(errors.New("error-creating-house"))
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
