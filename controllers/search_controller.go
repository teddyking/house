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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/teddyking/house/api/v1alpha1"
	housev1alpha1 "github.com/teddyking/house/api/v1alpha1"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Searcher
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SearchRepo

type Searcher interface {
	NumResults(url string) (int, error)
}

type SearchRepo interface {
	Get(ctx context.Context, namespacedName types.NamespacedName) (*v1alpha1.Search, error)
	UpdateStatus(ctx context.Context, search *v1alpha1.Search) error
}

// SearchReconciler reconciles a Search object
type SearchReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	Searcher   Searcher
	SearchRepo SearchRepo
}

// +kubebuilder:rbac:groups=house.teddyking.github.io,resources=searches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=house.teddyking.github.io,resources=searches/status,verbs=get;update;patch

func (r *SearchReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("search", req.NamespacedName)

	search, err := r.SearchRepo.Get(ctx, req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	numResults, err := r.Searcher.NumResults(search.Spec.URL)
	if err != nil {
		return ctrl.Result{}, err
	}

	search.Status.NumResults = numResults
	if err := r.SearchRepo.UpdateStatus(ctx, search); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SearchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&housev1alpha1.Search{}).
		Complete(r)
}
