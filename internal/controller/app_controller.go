/*
Copyright 2024 cara.

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
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1" // always needed
	corev1 "k8s.io/api/core/v1" // always needed

	v1 "github.com/testcara/app-operator/api/v1"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.wlin.cn,resources=apps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.wlin.cn,resources=apps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.wlin.cn,resources=apps/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;update;list;watch;create;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status, verbs=get
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;update;list;watch;create;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the App object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
var CounterReconcileApp int64

const GenericRequestDuration = 100 * time.Minute

func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	<-time.NewTicker(100 * time.Microsecond).C
	log := log.FromContext(ctx)

	CounterReconcileApp += 1
	log.Info("Starting a reconcile", "number", CounterReconcileApp)

	app := &v1.App{}
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		// when there is no App resource, return directly
		if errors.IsNotFound((err)) {
			log.Info("App not found")
			return ctrl.Result{}, nil
		}
		// when there is the App resource, retry every duration
		log.Error(err, "Failed to get the app, will request after a short time")
		return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
	}

	// reconcile sub-resources
	var result ctrl.Result
	var err error
	result, err = r.reconcileDeployment(ctx, app)
	if err != nil {
		log.Error(err, "Failed to reconcile Deployment")
		return result, err
	}
	result, err = r.reconcileService(ctx, app)
	if err != nil {
		log.Error(err, "Failed to reconcile Service")
		return result, err
	}

	log.Info("All resources have been reconciled")
	return ctrl.Result{}, nil
}

func (r *AppReconciler) reconcileDeployment(ctx context.Context, app *v1.App) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var dp = &appsv1.Deployment{}
	// go to the specified namespace to check whether the exact deployment exists
	err := r.Get(ctx, types.NamespacedName{
		Namespace: app.Namespace,
		Name:      app.Name,
	}, dp)
	// we get the deployment well
	if err == nil {
		log.Info("Deloyment has already exist")
		// if the current status == expected status, return directly
		if reflect.DeepEqual(dp.Status, app.Status.Workflow) {
			return ctrl.Result{}, nil
		}
		// when current status != expected status, update it
		app.Status.Workflow = dp.Status
		if err := r.Status().Update(ctx, app); err != nil {
			log.Error(err, "Failed to update App status")
			return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
		}
		log.Info("The App status has been updated")
		return ctrl.Result{}, nil
	}
	// if there is no deployment, we create one
	if errors.IsNotFound(err) {
		newDp := &appsv1.Deployment{}
		newDp.SetName(app.Name)
		newDp.SetNamespace(app.Namespace)
		newDp.SetLabels(app.Labels)
		newDp.Spec = app.Spec.Deployment.DeploymentSpec
		newDp.Spec.Template.SetLabels(app.Labels)

		if err := ctrl.SetControllerReference(app, newDp, r.Scheme); err != nil {
			log.Error(err, "Failed to SetControllerReference, will request after a short time")
			return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
		}
		if err := r.Create(ctx, newDp); err != nil {
			log.Error(err, "Failed to create Deployment, will request after a short time")
			return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
		}
		log.Info("Deployment has been created")
		return ctrl.Result{}, nil
	} else {
		// the deployment is there but we cannot get well
		log.Error(err, "Failed to get Deployment, will request after a short time")
		return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
	}

}

func (r *AppReconciler) reconcileService(ctx context.Context, app *v1.App) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var svc = &corev1.Service{}
	// go to the specified namespace to check whether the exact Svc exists
	err := r.Get(ctx, types.NamespacedName{
		Namespace: app.Namespace,
		Name:      app.Name,
	}, svc)
	// we get the service well
	if err == nil {
		log.Info("Service has already exist")
		// if the current status == expected status, return directly
		if reflect.DeepEqual(svc.Status, app.Status.Network) {
			return ctrl.Result{}, nil
		}
		// when current status != expected status, update it
		app.Status.Network = svc.Status
		if err := r.Status().Update(ctx, app); err != nil {
			log.Error(err, "Failed to update App status")
			return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
		}
		log.Info("The App status has been updated")
		return ctrl.Result{}, nil
	}
	// if there is no deployment, we create one
	if errors.IsNotFound(err) {
		newSvc := &corev1.Service{}
		newSvc.SetName(app.Name)
		newSvc.SetNamespace(app.Namespace)
		newSvc.SetLabels(app.Labels)
		newSvc.Spec = app.Spec.Service.ServiceSpec
		newSvc.Spec.Selector = app.Labels

		// the SetControllerReference is extrem important to ensure our control knows
		// which resources we need to manage
		if err := ctrl.SetControllerReference(app, newSvc, r.Scheme); err != nil {
			log.Error(err, "Failed to SetControllerReference, will request after a short time")
			return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
		}
		if err := r.Create(ctx, newSvc); err != nil {
			log.Error(err, "Failed to create Service, will request after a short time")
			return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
		}
		log.Info("Service has been created")
		return ctrl.Result{}, nil
	} else {
		// the deployment is there but we cannot get well
		log.Error(err, "Failed to get Service, will request after a short time")
		return ctrl.Result{RequeueAfter: GenericRequestDuration}, err
	}

}

// SetupWithManager sets up the controller with the Manager.
func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.App{}).
		Complete(r)
}
