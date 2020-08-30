package repository

import (
	"context"

	"github.com/teddyking/house/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Search struct {
	client.Client
}

func (s *Search) Get(ctx context.Context, namespacedName types.NamespacedName) (*v1alpha1.Search, error) {
	search := &v1alpha1.Search{}

	if err := s.Client.Get(ctx, namespacedName, search); err != nil {
		return &v1alpha1.Search{}, err
	}

	return search, nil
}

func (s *Search) UpdateStatus(ctx context.Context, search *v1alpha1.Search) error {
	return s.Client.Status().Update(ctx, search)
}
