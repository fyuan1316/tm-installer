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

package v1alpha1

import (
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type TaskSpec struct {
	Resources   []string   `json:"resources,omitempty"`
	JobTemplate v1.JobSpec `json:"job_template"`
}

// TaskRunSpec defines the desired state of TaskRun
type TaskRunSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Steps []TaskSpec `json:"steps"`
}

//type Condition apis.Condition

// TaskRunStatus defines the observed state of TaskRun
type TaskRunStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions the latest available observations of a resource's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Condition           *Condition   `json:"condition,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	State               TaskRunState `json:"state"`
	TaskRunStatusFields `json:",inline"`
}
type TaskRunState string

const (
	Installing TaskRunState = "Installing"
	Succeed    TaskRunState = "Succeed"
	Failed     TaskRunState = "Failed"
)

type TaskRunStatusFields struct {
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is the time the build completed.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Steps describes the state of each build step container.
	// +optional
	Steps []StepState `json:"steps,omitempty"`
}
type StepState struct {
	corev1.ContainerState `json:",inline"`
	Condition             Condition `json:"condition,omitempty"`

	Name          string `json:"name,omitempty"`
	ContainerName string `json:"container,omitempty"`
	ImageID       string `json:"imageID,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// TaskRun is the Schema for the taskruns API
type TaskRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskRunSpec   `json:"spec,omitempty"`
	Status TaskRunStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskRunList contains a list of TaskRun
type TaskRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TaskRun{}, &TaskRunList{})
}
