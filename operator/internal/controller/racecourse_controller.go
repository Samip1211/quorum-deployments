/*
Copyright 2023.

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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/Samip1211/racecourse/api/v1alpha1"
	kaleidov1alpha1 "github.com/Samip1211/racecourse/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

var namespace string

// RacecourseReconciler reconciles a Racecourse object
type RacecourseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kaleido.kaleido.com,resources=racecourses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kaleido.kaleido.com,resources=racecourses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kaleido.kaleido.com,resources=racecourses/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Racecourse object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *RacecourseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	raceObj := &v1alpha1.Racecourse{}
	err := r.Client.Get(ctx, req.NamespacedName, raceObj)
	if err != nil {
		return ctrl.Result{}, err
	}

	namespace = req.Namespace

	// if raceObj.Status.DeploymentStatus == "True" {
	// 	return ctrl.Result{}, nil
	// }

	// Reconcile the Deployment.
	if raceObj.Status.DeploymentStatus == "" {
		err = r.reconcileRaceCourseDeployment(*raceObj)
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}
		// Reconcile Service
		r.reconcileRaceCourseService()

		raceObj.Status.DeploymentStatus = "Waiting"
		err = r.Status().Update(ctx, raceObj)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Get the Deployment.
	deploy := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "racecourse"}, deploy)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check the status of the deployment.
	if (deploy.Status.Replicas == deploy.Status.ReadyReplicas) && deploy.Status.UnavailableReplicas == 0 {
		raceObj.Status.DeploymentStatus = "True"
		return ctrl.Result{Requeue: false}, r.Status().Update(ctx, raceObj)
	}

	return ctrl.Result{Requeue: true}, nil

}
func (r *RacecourseReconciler) reconcileRaceCourseService() {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "racecourse",
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": "racecourse"},
			Ports: []v1.ServicePort{
				{
					Port:       3000,
					TargetPort: intstr.FromInt(3000),
				},
			},
		},
	}
	r.Create(context.Background(), service)
}

func (r *RacecourseReconciler) reconcileRaceCourseDeployment(rcourse v1alpha1.Racecourse) error {
	replicasNum := int32(rcourse.Spec.Replicas)
	deploy := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "racecourse",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicasNum,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "racecourse"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "racecourse"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "racecourse",
							Image:           rcourse.Spec.Image,
							ImagePullPolicy: v1.PullIfNotPresent,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 3000,
									Name:          "http",
								},
							},
						},
					},
				},
			},
		},
	}

	return r.Client.Create(context.Background(), &deploy)
}

// SetupWithManager sets up the controller with the Manager.
func (r *RacecourseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kaleidov1alpha1.Racecourse{}).
		Complete(r)
}
