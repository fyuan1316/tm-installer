package v1alpha1

import corev1 "k8s.io/api/core/v1"

func (tr TaskRun) initCondition() {
	if tr.Status.Steps == nil {
		for _, spec := range tr.Spec.Steps {
			cond := Condition{}
			cond.Type = ConditionSucceeded
			cond.Status = corev1.ConditionUnknown
			stepState := StepState{}
			stepState.Condition = cond
			stepState.Name = spec.JobTemplate.Template.Name
			tr.Status.Steps = append(tr.Status.Steps, stepState)
		}
	}
	if tr.Status.Condition == nil {
		tr.Status.Condition = &Condition{}
		tr.Status.Condition.Type = ConditionSucceeded
		tr.Status.Condition.Status = corev1.ConditionUnknown
	}
}
func (trs TaskRunStatus) IsInstalling() bool {
	return trs.State == Installing
}

func (trs TaskRunStatus) IsAllStepsDone() bool {
	done := true
	for _, stepStatus := range trs.Steps {
		if stepStatus.Condition.Type != ConditionSucceeded ||
			stepStatus.Condition.Status == corev1.ConditionUnknown ||
			stepStatus.Condition.Status == ConditionScheduled {
			done = false
		}
	}
	return done
}
func (trs *TaskRunStatus) GetStep(task TaskSpec) *StepState {
	for i, step := range trs.Steps {
		if step.Name == task.JobTemplate.Template.Name {
			return &trs.Steps[i]
		}
	}
	return nil
}
func (trss StepState) IsNotScheduled() bool {
	if trss.Condition.Type == ConditionSucceeded && trss.Condition.Status == corev1.ConditionUnknown {
		return true
	}
	return false
}
func (trss StepState) IsScheduled() bool {
	if trss.Condition.Type == ConditionSucceeded && trss.Condition.Status == ConditionScheduled {
		return true
	}
	return false
}
