/*
Copyright 2022.

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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	interviewv1alpha1 "github.com/anynines/tmp-homework-ms/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DummyReconciler reconciles a Dummy object
type DummyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const (
	DefaultRequeueTime = 20 * time.Second
	DefaultRequeue     = true
	Finalizer          = "interview.anynines.com/dummy-finalizer"
)

func increment() int {
	value := 0
	value++
	return value
}

//+kubebuilder:rbac:groups=interview.interview.com,resources=dummies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=interview.interview.com,resources=dummies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=interview.interview.com,resources=dummies/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DummyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Dummy", req.NamespacedName)
	log.Info("Reconciling Dummy")

	dummyInstance := &interviewv1alpha1.Dummy{}
	if err := r.Get(ctx, req.NamespacedName, dummyInstance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{Requeue: DefaultRequeue, RequeueAfter: DefaultRequeueTime}, nil
		}
		log.Error(err, fmt.Sprintf("Failed to get Dummy: %s", req.Name))
		return ctrl.Result{Requeue: DefaultRequeue, RequeueAfter: DefaultRequeueTime}, err
	}

	log.Info(fmt.Sprintf("Found resource with name %s in namespace %s", dummyInstance.Name, dummyInstance.Namespace))
	if dummyInstance.Spec.Message != nil {
		log.Info(fmt.Sprintf("Found message %s", *dummyInstance.Spec.Message))
		dummyPatchBase := client.MergeFrom(dummyInstance.DeepCopy())
		dummyInstance.Status.SpecEcho = *dummyInstance.Spec.Message
		if err := r.Client.Status().Patch(ctx, dummyInstance, dummyPatchBase); err != nil {
			log.Error(err, fmt.Sprintf("Failed to update status for Dummy: %s", dummyInstance.Name))
			return ctrl.Result{Requeue: DefaultRequeue, RequeueAfter: DefaultRequeueTime}, err
		}
	}

	pod := &corev1.Pod{}

	if err := r.Get(ctx, req.NamespacedName, pod); err != nil {
		if errors.IsNotFound(err) {
			if err := r.createPod(ctx, log, dummyInstance); err != nil {
				log.Error(err, fmt.Sprintf("Failed to create Pod: %s", req.Name))
				return ctrl.Result{Requeue: DefaultRequeue, RequeueAfter: DefaultRequeueTime}, err
			}
			return ctrl.Result{Requeue: DefaultRequeue, RequeueAfter: DefaultRequeueTime}, nil
		}
		log.Error(err, fmt.Sprintf("Failed to get Dummy: %s", req.Name))
		return ctrl.Result{Requeue: DefaultRequeue, RequeueAfter: DefaultRequeueTime}, err
	}

	r.UpdateDummyPodStatus(ctx, log, dummyInstance, string(pod.Status.Phase))

	return ctrl.Result{Requeue: DefaultRequeue, RequeueAfter: DefaultRequeueTime}, nil
}

func (r *DummyReconciler) createPod(ctx context.Context, log logr.Logger, dummy *interviewv1alpha1.Dummy) error {

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dummy.Name,
			Namespace: dummy.Namespace,
			Labels: map[string]string{
				"run": "nginx",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx-container",
					Image: "nginx",
				},
			},
			DNSPolicy:     corev1.DNSClusterFirst,
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}
	controllerutil.SetOwnerReference(dummy, pod, r.Scheme)
	if err := r.Client.Create(ctx, pod); err != nil {
		return err
	}
	if err := r.UpdateDummyPodStatus(ctx, log, dummy, string(pod.Status.Phase)); err != nil {
		return err
	}
	return nil
}

func (r *DummyReconciler) UpdateDummyPodStatus(ctx context.Context, log logr.Logger, dummy *interviewv1alpha1.Dummy, podStatus string) error {
	dummyPatchBase := client.MergeFrom(dummy.DeepCopy())
	dummy.Status.PodStatus = podStatus
	if err := r.Client.Status().Patch(ctx, dummy, dummyPatchBase); err != nil {
		log.Error(err, fmt.Sprintf("Failed to update status for Dummy: %s", dummy.Name))
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&interviewv1alpha1.Dummy{}).
		Complete(r)
}
