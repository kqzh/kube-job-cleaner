/*
Copyright 2021 kqzh.

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

package utils

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsJobFailed(j *batchv1.Job) bool {
	for _, c := range j.Status.Conditions {
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func IsJobSuccess(j *batchv1.Job) bool {
	for _, c := range j.Status.Conditions {
		if c.Type == batchv1.JobComplete && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func IsJobFinished(j *batchv1.Job) bool {
	if IsJobSuccess(j) || IsJobFailed(j) {
		return true
	}
	return false
}

func GetConfigMaps(j *batchv1.Job) []*corev1.ConfigMap {
	configMaps := make([]*corev1.ConfigMap, 0)
	for _, volume := range j.Spec.Template.Spec.Volumes {
		if volume.ConfigMap != nil {
			configMaps = append(configMaps, &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      volume.ConfigMap.Name,
					Namespace: j.Namespace,
				},
			})
		}
	}
	return configMaps
}
