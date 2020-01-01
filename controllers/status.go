package controllers

import (
	installv1alpha1 "github.com/fyuan1316/tm-installer/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func initCondition(ctx *TaskRunContext, tr *installv1alpha1.TaskRun) {
	if tr.Status.Steps == nil {
		for _, spec := range tr.Spec.Steps {
			cond := installv1alpha1.Condition{}
			cond.Type = installv1alpha1.ConditionSucceeded
			cond.Status = corev1.ConditionUnknown
			cond.LastTransitionTime = metav1.Now()
			stepState := installv1alpha1.StepState{}
			stepState.Condition = cond
			stepState.Name = spec.JobTemplate.Template.Name
			tr.Status.Steps = append(tr.Status.Steps, stepState)
		}
		ctx.update = true
	}
	if tr.Status.Condition == nil {
		tr.Status.Condition = &installv1alpha1.Condition{}
		tr.Status.Condition.Type = installv1alpha1.ConditionSucceeded
		tr.Status.Condition.Status = corev1.ConditionUnknown
		tr.Status.Condition.LastTransitionTime = metav1.Now()
		ctx.update = true
	}
}

func convertToDone(ctx *TaskRunContext, trStatus *installv1alpha1.TaskRunStatus) {
	succeed := true
	for _, stepState := range trStatus.Steps {
		if stepState.Condition.Status == corev1.ConditionFalse {
			succeed = false
			break
		}
	}
	trStatus.Condition.Status = corev1.ConditionTrue
	if succeed {
		trStatus.State = installv1alpha1.Succeed
	} else {
		trStatus.State = installv1alpha1.Failed
	}
	ctx.update = true
}

func convertStepToScheduled(ctx *TaskRunContext, stepState *installv1alpha1.StepState) {
	stepState.Condition.Status = installv1alpha1.ConditionScheduled
	ctx.update = true
}
func convertStepToSucceed(ctx *TaskRunContext, stepState *installv1alpha1.StepState) {
	stepState.Condition.Status = corev1.ConditionTrue
	ctx.update = true
}
func convertStepToFailed(ctx *TaskRunContext, stepState *installv1alpha1.StepState) {
	stepState.Condition.Status = corev1.ConditionFalse
	ctx.update = true
}
