package v1alpha1

import corev1 "k8s.io/api/core/v1"

func (tr TaskRun) GetNextTask() *TaskSpec {
	if tr.Status.Steps == nil {
		return nil
	}
	for i, t := range tr.Status.Steps {
		if t.Condition.Type == ConditionSucceeded &&
			t.Condition.Status != corev1.ConditionTrue &&
			t.Condition.Status != corev1.ConditionFalse {
			return &tr.Spec.Steps[i]
		}
	}
	return nil
}
