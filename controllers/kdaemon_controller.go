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

package controllers

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/kibaamor/operator-demo/api/v1alpha1"
)

// KDaemonReconciler reconciles a KDaemon object
type KDaemonReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.kibazen.cn,resources=kdaemons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.kibazen.cn,resources=kdaemons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.kibazen.cn,resources=kdaemons/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KDaemon object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *KDaemonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	kds := &appsv1alpha1.KDaemon{}
	if err := r.Client.Get(ctx, req.NamespacedName, kds); err != nil {
		log.Error(err, "failed to get KDaemon")
		return ctrl.Result{
			Requeue: true,
		}, err
	}
	if kds.Spec.Image == "" {
		err := fmt.Errorf("invalid image config")
		log.Error(err, "can not deploy with empty image")
		return ctrl.Result{}, err
	}

	nl := &v1.NodeList{}
	if err := r.Client.List(ctx, nl); err != nil {
		log.Error(err, "failed to get node list")
		return ctrl.Result{
			Requeue: true,
		}, err
	}

	for _, n := range nl.Items {
		p := &v1.Pod{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("%s-", n.Name),
				Namespace:    kds.Namespace,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Image: kds.Spec.Image,
						Name:  "kdaemon",
					},
				},
				NodeName: n.Name,
			},
		}
		if err := r.Client.Create(ctx, p); err != nil {
			log.Error(err, "failed create pod on Node", "node", n.Name)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KDaemonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.KDaemon{}).
		Complete(r)
}
