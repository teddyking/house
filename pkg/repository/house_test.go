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
	Describe("Create", func() {
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

			Expect(repo.Create(context.TODO(), house)).To(Succeed())

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
			BeforeEach(func() {
				house := &v1alpha1.House{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "house-1",
						Namespace: "default",
					},
					Spec:   v1alpha1.HouseSpec{},
					Status: v1alpha1.HouseStatus{},
				}

				Expect(fakeClient.Create(context.TODO(), house)).To(Succeed())
			})

			It("does nothing", func() {
				Expect(repo.Create(context.TODO(), htypes.House{ID: "house-1"})).To(Succeed())
			})
		})
	})
})
