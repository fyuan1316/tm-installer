package controllers

import (
	"context"
	installv1alpha1 "github.com/fyuan1316/tm-installer/api/v1alpha1"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type TaskRunContext struct {
	update         bool
	context        context.Context
	log            logr.Logger
	request        ctrl.Request
	taskRunPointer *installv1alpha1.TaskRun
}

func newContext(req ctrl.Request, log logr.Logger) *TaskRunContext {
	return &TaskRunContext{
		update:         false,
		context:        context.Background(),
		log:            log,
		request:        req,
		taskRunPointer: &installv1alpha1.TaskRun{},
	}
}

func (trctx *TaskRunContext) TaskRunCopy() *installv1alpha1.TaskRun {
	return trctx.taskRunPointer.DeepCopy()
}
