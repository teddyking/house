package repository_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/teddyking/house/api/v1alpha1"
	. "github.com/teddyking/house/pkg/repository"
	htypes "github.com/teddyking/house/pkg/types"
)

var _ = Describe("House", func() {
	Describe("Upsert", func() {
		var (
			fakeClient client.Client
			repo       *House
		)

		BeforeEach(func() {
			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			fakeClient = fake.NewFakeClientWithScheme(scheme)

			repo = &House{Client: fakeClient}
		})

		It("creates the House in the k8s API", func() {
			house := htypes.House{
				ID:          "house-1",
				Price:       "price-1",
				OfferType:   "offer-type-1",
				Description: "description-1",
				Postcode:    "postcode-1",
				URL:         "url-1",
			}

			Expect(repo.Upsert(context.TODO(), house)).To(Succeed())

			createdHouse := &v1alpha1.House{}
			Expect(fakeClient.Get(
				context.TODO(),
				types.NamespacedName{Name: "house-1", Namespace: "default"},
				createdHouse,
			)).To(Succeed())

			Expect(createdHouse.Name).To(Equal("house-1"))
			Expect(createdHouse.Spec.Price).To(Equal("price-1"))
			Expect(createdHouse.Spec.OfferType).To(Equal("offer-type-1"))
			Expect(createdHouse.Spec.Description).To(Equal("description-1"))
			Expect(createdHouse.Spec.Postcode).To(Equal("postcode-1"))
			Expect(createdHouse.Spec.URL).To(Equal("url-1"))
		})

		When("the House already exists in the k8s API", func() {
			var house *v1alpha1.House

			BeforeEach(func() {
				house = &v1alpha1.House{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "house-1",
						Namespace: "default",
					},
					Spec: v1alpha1.HouseSpec{
						Price: "price-1",
					},
					Status: v1alpha1.HouseStatus{},
				}

				Expect(fakeClient.Create(context.TODO(), house)).To(Succeed())
			})

			It("updates the House in the k8s API", func() {
				Expect(repo.Upsert(context.TODO(),
					htypes.House{
						ID:          "house-1",
						Price:       "price-2",
						OfferType:   "offer-type-2",
						Description: "description-2",
						Postcode:    "postcode-2",
						URL:         "url-2",
					})).To(Succeed())

				updatedHouse := &v1alpha1.House{}
				Expect(fakeClient.Get(
					context.TODO(),
					types.NamespacedName{Name: "house-1", Namespace: "default"},
					updatedHouse,
				)).To(Succeed())

				Expect(updatedHouse.Name).To(Equal("house-1"))
				Expect(updatedHouse.Spec.Price).To(Equal("price-2"))
				Expect(updatedHouse.Spec.OfferType).To(Equal("offer-type-2"))
				Expect(updatedHouse.Spec.Description).To(Equal("description-2"))
				Expect(updatedHouse.Spec.Postcode).To(Equal("postcode-2"))
				Expect(updatedHouse.Spec.URL).To(Equal("url-2"))
			})
		})
	})
})
