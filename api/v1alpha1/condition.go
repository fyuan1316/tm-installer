package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType is a camel-cased condition type.
type ConditionType string

const (
	// ConditionReady specifies that the resource is ready.
	// For long-running resources.
	ConditionReady ConditionType = "Ready"
	// ConditionSucceeded specifies that the resource has finished.
	// For resource which run to completion.
	ConditionSucceeded ConditionType = "Succeeded"

	ConditionScheduled corev1.ConditionStatus = "Scheduled"
)

// ConditionSeverity expresses the severity of a Condition Type failing.
type ConditionSeverity string

const (
	// ConditionSeverityError specifies that a failure of a condition type
	// should be viewed as an error.  As "Error" is the default for conditions
	// we use the empty string (coupled with omitempty) to avoid confusion in
	// the case where the condition is in state "True" (aka nothing is wrong).
	ConditionSeverityError ConditionSeverity = ""
	// ConditionSeverityWarning specifies that a failure of a condition type
	// should be viewed as a warning, but that things could still work.
	ConditionSeverityWarning ConditionSeverity = "Warning"
	// ConditionSeverityInfo specifies that a failure of a condition type
	// should be viewed as purely informational, and that things could still work.
	ConditionSeverityInfo ConditionSeverity = "Info"
)

// +k8s:deepcopy-gen=true
type Condition struct {
	// Type of condition.
	// +required
	Type ConditionType `json:"type" description:"type of status condition"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status" description:"status of the condition, one of True, False, Unknown"`

	// Severity with which to treat failures of this type of condition.
	// When this is not specified, it defaults to Error.
	// +optional
	Severity ConditionSeverity `json:"severity,omitempty" description:"how to interpret failures of this condition, one of Error, Warning, Info"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// We use VolatileTime in place of metav1.Time to exclude this from creating equality.Semantic
	// differences (all other things held constant).
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" description:"last time the condition transit from one status to another"`

	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" description:"one-word CamelCase reason for the condition's last transition"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" description:"human-readable message indicating details about last transition"`
}
