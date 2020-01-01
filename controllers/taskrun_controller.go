/*
Copyright 2020 fyuan.

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
	"github.com/prometheus/common/log"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	installv1alpha1 "github.com/fyuan1316/tm-installer/api/v1alpha1"
	"golang.org/x/xerrors"
)

// TaskRunReconciler reconciles a TaskRun object
type TaskRunReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=install.alauda.io,resources=taskruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=install.alauda.io,resources=taskruns/status,verbs=get;update;patch

func (r *TaskRunReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	trLog := r.Log.WithValues("taskRun", req.NamespacedName)
	ctx := newContext(req, trLog)

	// your logic here
	if err := r.Client.Get(ctx.context, req.NamespacedName, ctx.taskRunPointer); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	trCopy := ctx.TaskRunCopy()
	initCondition(ctx, trCopy)

	if trCopy.Status.IsInstalling() {
		if trCopy.Status.IsAllStepsDone() {
			convertToDone(ctx, &trCopy.Status)
		} else {
			task := trCopy.GetNextTask()
			if err := r.reconcileStepTask(ctx, trCopy, *task); err != nil {
				log.Error(err, "reconcileStepTask err")
				return ctrl.Result{}, err
			}
		}
	}
	if ctx.update {
		if err := r.Update(ctx.context, trCopy); err != nil {
			log.Error(err, "update taskRun error")
		}
	}

	return ctrl.Result{}, nil
}
func (r *TaskRunReconciler) reconcileStepTask(ctx *TaskRunContext, tr *installv1alpha1.TaskRun, task installv1alpha1.TaskSpec) error {
	if task.Resources != nil {
		for i, res := range task.Resources {
			unStruct := &unstructured.Unstructured{}
			if err := yaml.Unmarshal([]byte(res), unStruct); err != nil {
				log.Error(err, "cast resource%d of %s error", i, task.JobTemplate.Template.Name)
			}
			controllerutil.SetControllerReference(tr, unStruct, r.Scheme)
			if err := r.Client.Create(context.Background(), unStruct); errors.IsAlreadyExists(err) {
				if err := r.Client.Update(context.Background(), unStruct); err != nil {
					log.Error(err, "update resource%d of %s error", i, task.JobTemplate.Template.Name)
				}
			}
		}
	}
	step := tr.Status.GetStep(task)
	if step == nil {
		return xerrors.Errorf("stepState of %s is nil", task.JobTemplate.Template.Name)
	}
	if step.IsNotScheduled() {
		job := renderJob(task)
		controllerutil.SetControllerReference(tr, job, r.Scheme)
		if err := r.Client.Create(ctx.context, job); err != nil {
			log.Error(err, "create job %s error", task.JobTemplate.Template.Name)
		}
		convertStepToScheduled(ctx, step)
	}
	if step.IsScheduled() {
		job := &v1.Job{}
		template := task.JobTemplate.Template
		objectKey := client.ObjectKey{Name: getName(template.Name), Namespace: template.Namespace}
		err := r.Client.Get(ctx.context, objectKey, job)
		if errors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return xerrors.Errorf("get job error when process scheduled step: %w", err)
		}
		if job.Status.Succeeded > 0 {
			convertStepToSucceed(ctx, step)
		}
		if job.Status.Failed > 0 {
			convertStepToFailed(ctx, step)
		}
	}

	return nil
}

func (r *TaskRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&installv1alpha1.TaskRun{}).
		Owns(&v1.Job{}).
		Complete(r)
}
