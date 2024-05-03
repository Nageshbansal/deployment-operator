/*
Copyright 2024.

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

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	deployv1 "github.com/Nageshbansal/deployment-operator/api/v1"
)

// DeploySetReconciler reconciles a DeploySet object
type DeploySetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	DefaultReconciliationInterval = 5
)

//+kubebuilder:rbac:groups=deploy.nagesh-node.me,resources=deploysets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=deploy.nagesh-node.me,resources=deploysets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=deploy.nagesh-node.me,resources=deploysets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeploySet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *DeploySetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("starting the reconcilation")

	ds := &deployv1.DeploySet{}

	err := r.getDeploySet(ctx, req, ds)
	if err != nil {
		return ctrl.Result{}, err
	}

	ok, err := r.DeploymentIfNotExist(ctx, req, ds)

	if err != nil {
		log.Error(err, "failed to deploy the deployment for deploySet")
		return ctrl.Result{}, err
	}

	if ok {
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	err = r.UpdateDeploymentReplica(ctx, req, ds)
	if err != nil {
		log.Error(err, "failed to update the deployment for deploySet")

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploySetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&deployv1.DeploySet{}).
		Complete(r)
}

func (r *DeploySetReconciler) getDeploySet(ctx context.Context, req ctrl.Request, ds *deployv1.DeploySet) error {

	err := r.Get(ctx, req.NamespacedName, ds)
	if err != nil {
		return err
	}
	return nil
}
