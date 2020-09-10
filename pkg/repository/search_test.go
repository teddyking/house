package repository_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/teddyking/house/api/v1alpha1"
	. "github.com/teddyking/house/pkg/repository"
)

var _ = Describe("Search", func() {
	Describe("Get", func() {
		var (
			search *v1alpha1.Search

			namespacedName types.NamespacedName

			fakeClient     client.Client
			returnedErr    error
			returnedSearch *v1alpha1.Search

			repo *Search
		)

		BeforeEach(func() {
			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			name := "search-1"
			namespace := "default"
			namespacedName = types.NamespacedName{Name: name, Namespace: namespace}

			search = &v1alpha1.Search{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: v1alpha1.SearchSpec{
					URL: "url-1",
				},
				Status: v1alpha1.SearchStatus{},
			}

			fakeClient = fake.NewFakeClientWithScheme(scheme)
			Expect(fakeClient.Create(context.TODO(), search)).To(Succeed())

			repo = &Search{Client: fakeClient}
		})

		JustBeforeEach(func() {
			returnedSearch, returnedErr = repo.Get(context.TODO(), namespacedName)
		})

		When("the Search exists", func() {
			It("gets the Search from the k8s API", func() {
				Expect(returnedErr).NotTo(HaveOccurred())

				Expect(returnedSearch.Name).To(Equal(search.Name))
				Expect(returnedSearch.Spec.URL).To(Equal(search.Spec.URL))
			})
		})

		When("the Search doesn't exist", func() {
			BeforeEach(func() {
				Expect(fakeClient.Delete(context.TODO(), search)).To(Succeed())
			})

			It("returns a NotFound error", func() {
				Expect(returnedErr).To(HaveOccurred())
				Expect(apierrors.IsNotFound(returnedErr)).To(BeTrue())
			})
		})
	})

	Describe("UpdateStatus", func() {
		var (
			search         *v1alpha1.Search
			namespacedName types.NamespacedName

			fakeClient client.Client
			repo       *Search
		)

		BeforeEach(func() {
			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)

			name := "search-1"
			namespace := "default"
			namespacedName = types.NamespacedName{Name: name, Namespace: namespace}

			search = &v1alpha1.Search{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec:   v1alpha1.SearchSpec{},
				Status: v1alpha1.SearchStatus{},
			}

			fakeClient = fake.NewFakeClientWithScheme(scheme)
			Expect(fakeClient.Create(context.TODO(), search)).To(Succeed())

			repo = &Search{Client: fakeClient}
		})

		It("updates the Search status in the k8s API", func() {
			search.Status.NumResults = 5
			search.Status.ObservedGeneration = 2
			Expect(repo.UpdateStatus(context.TODO(), search)).To(Succeed())

			updatedSearch := &v1alpha1.Search{}
			Expect(fakeClient.Get(context.TODO(), namespacedName, updatedSearch)).To(Succeed())

			Expect(updatedSearch.Status.NumResults).To(Equal(5))
			Expect(updatedSearch.Status.ObservedGeneration).To(BeEquivalentTo(2))
		})
	})
})
