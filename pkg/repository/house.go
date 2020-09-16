package repository

import (
	"context"

	"github.com/teddyking/house/api/v1alpha1"
	htypes "github.com/teddyking/house/pkg/types"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type House struct {
	client.Client
}

func (h *House) Upsert(ctx context.Context, house htypes.House) error {
	currentHouse := &v1alpha1.House{}
	err := h.Get(ctx, types.NamespacedName{Name: house.ID, Namespace: "default"}, currentHouse)

	// CREATE
	if apierrors.IsNotFound(err) {
		return h.Client.Create(ctx, &v1alpha1.House{
			ObjectMeta: metav1.ObjectMeta{
				Name:      house.ID,
				Namespace: "default",
			},
			Spec: v1alpha1.HouseSpec{
				Price:       house.Price,
				OfferType:   house.OfferType,
				Description: house.Description,
				Postcode:    house.Postcode,
				URL:         house.URL,
			},
		})
	}

	// UPDATE
	currentHouse.Spec.Price = house.Price
	currentHouse.Spec.OfferType = house.OfferType
	currentHouse.Spec.Description = house.Description
	currentHouse.Spec.Postcode = house.Postcode
	currentHouse.Spec.URL = house.URL

	return h.Client.Update(ctx, currentHouse)
}
