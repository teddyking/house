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

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/teddyking/house/api/v1alpha1"
	housev1alpha1 "github.com/teddyking/house/api/v1alpha1"
	htypes "github.com/teddyking/house/pkg/types"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Scraper
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SearchRepo
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . HouseRepo

type Scraper interface {
	Properties(url string) ([]htypes.House, error)
}

type SearchRepo interface {
	Get(ctx context.Context, namespacedName types.NamespacedName) (*v1alpha1.Search, error)
	UpdateStatus(ctx context.Context, search *v1alpha1.Search) error
}

type HouseRepo interface {
	Upsert(ctx context.Context, house htypes.House) error
}

// SearchReconciler reconciles a Search object
type SearchReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	Scraper Scraper

	SearchRepo SearchRepo
	HouseRepo  HouseRepo
}

// +kubebuilder:rbac:groups=house.teddyking.github.io,resources=searches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=house.teddyking.github.io,resources=houses,verbs=create;update
// +kubebuilder:rbac:groups=house.teddyking.github.io,resources=searches/status,verbs=get;update;patch

func (r *SearchReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	r.Log = r.Log.WithValues("search", req.NamespacedName)

	r.Log.Info("Reconcile - start")

	search, err := r.SearchRepo.Get(ctx, req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	properties, err := r.Scraper.Properties(search.Spec.URL)
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, house := range properties {
		if err := r.HouseRepo.Upsert(ctx, house); err != nil {
			return ctrl.Result{}, err
		}
	}

	search.Status.NumResults = len(properties)
	search.Status.ObservedGeneration = search.Generation

	if err := r.SearchRepo.UpdateStatus(ctx, search); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("Reconcile - end")

	return ctrl.Result{RequeueAfter: time.Hour}, nil
}

func (r *SearchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&housev1alpha1.Search{}).
		Complete(r)
}
