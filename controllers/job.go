package controllers

import (
	"github.com/fyuan1316/tm-installer/api/v1alpha1"
	v1 "k8s.io/api/batch/v1"
	"strings"
)

func renderJob(task v1alpha1.TaskSpec) *v1.Job {
	job := &v1.Job{}
	job.Spec = task.JobTemplate
	job.Name = getName(task.JobTemplate.Template.Name)
	job.Namespace = task.JobTemplate.Template.Namespace
	return job
}
func getName(name string) string {
	name = strings.ToLower(name)
	return name
}
