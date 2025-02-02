/*
Copyright 2022 The KServe Authors.

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
package pod

import (
	"github.com/kserve/kserve/pkg/constants"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmp"

	"testing"
)

const sklearnPrometheusPort = "8080"

func TestInjectMetricsAggregator(t *testing.T) {
	scenarios := map[string]struct {
		original *v1.Pod
		expected *v1.Pod
	}{
		"EnableMetricAggTrue": {
			original: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation: "true",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
						},
					},
				},
			},
			expected: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation: "true",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
							Env: []v1.EnvVar{
								{Name: constants.KServeContainerPrometheusMetricsPortEnvVarKey, Value: sklearnPrometheusPort},
								{Name: constants.KServeContainerPrometheusMetricsPathEnvVarKey, Value: constants.DefaultPrometheusPath},
								{Name: constants.QueueProxyAggregatePrometheusMetricsPortEnvVarKey, Value: constants.QueueProxyAggregatePrometheusMetricsPort},
							},
						},
					},
				},
			},
		},
		"EnableMetricAggNotSet": {
			original: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
						},
					},
				},
			},
			expected: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation: "true",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
						},
					},
				},
			},
		},
		"EnableMetricAggFalse": {
			original: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation: "false",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
						},
					},
				},
			},
			expected: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation: "true",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
						},
					},
				},
			},
		},
		"setPromAnnotationTrue": {
			original: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation:          "true",
						constants.SetPrometheusAggregateAnnotation: "true",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
						},
					},
				},
			},
			expected: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation:          "true",
						constants.SetPrometheusAggregateAnnotation: "true",
						constants.PrometheusPortAnnotationKey:      constants.QueueProxyAggregatePrometheusMetricsPort,
						constants.PrometheusPathAnnotationKey:      constants.DefaultPrometheusPath,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
							Env: []v1.EnvVar{
								{Name: constants.KServeContainerPrometheusMetricsPortEnvVarKey, Value: sklearnPrometheusPort},
								{Name: constants.KServeContainerPrometheusMetricsPathEnvVarKey, Value: constants.DefaultPrometheusPath},
								{Name: constants.QueueProxyAggregatePrometheusMetricsPortEnvVarKey, Value: constants.QueueProxyAggregatePrometheusMetricsPort},
							},
						},
					},
				},
			},
		},
		"SetPromAnnotationFalse": {
			original: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation:          "true",
						constants.SetPrometheusAggregateAnnotation: "false",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
						},
					},
				},
			},
			expected: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment",
					Namespace: "default",
					Annotations: map[string]string{
						constants.EnableMetricAggregation:          "true",
						constants.SetPrometheusAggregateAnnotation: "false",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name: "sklearn",
					},
						{
							Name: "queue-proxy",
							Env: []v1.EnvVar{
								{Name: constants.KServeContainerPrometheusMetricsPortEnvVarKey, Value: sklearnPrometheusPort},
								{Name: constants.KServeContainerPrometheusMetricsPathEnvVarKey, Value: constants.DefaultPrometheusPath},
								{Name: constants.QueueProxyAggregatePrometheusMetricsPortEnvVarKey, Value: constants.QueueProxyAggregatePrometheusMetricsPort},
							},
						},
					},
				},
			},
		},
	}

	for name, scenario := range scenarios {
		err := InjectMetricsAggregator(scenario.original)
		if err != nil {
			t.Errorf("Test %q unexpected error %e", name, err)
		}
		if diff, _ := kmp.SafeDiff(scenario.expected.Spec, scenario.original.Spec); diff != "" {
			t.Errorf("Test %q unexpected result (-want +got): %v", name, diff)
		}
	}
}
