/*
Copyright 2023 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    htcp://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
	testutils "sigs.k8s.io/jobset/pkg/util/testing"
)

func TestIsJobFinished(t *testing.T) {
	tests := []struct {
		name              string
		conditions        []batchv1.JobCondition
		finished          bool
		wantConditionType batchv1.JobConditionType
	}{
		{
			name: "succeeded",
			conditions: []batchv1.JobCondition{
				{
					Type:   batchv1.JobComplete,
					Status: corev1.ConditionTrue,
				},
			},
			finished:          true,
			wantConditionType: batchv1.JobComplete,
		},
		{
			name: "failed",
			conditions: []batchv1.JobCondition{
				{
					Type:   batchv1.JobFailed,
					Status: corev1.ConditionTrue,
				},
			},
			finished:          true,
			wantConditionType: batchv1.JobFailed,
		},
		{
			name: "active",
			conditions: []batchv1.JobCondition{
				{
					Type:   "",
					Status: corev1.ConditionTrue,
				},
			},
			finished:          false,
			wantConditionType: "",
		},
		{
			name: "suspended",
			conditions: []batchv1.JobCondition{
				{
					Type:   batchv1.JobSuspended,
					Status: corev1.ConditionTrue,
				},
			},
			finished:          false,
			wantConditionType: "",
		},
		{
			name: "failure target",
			conditions: []batchv1.JobCondition{
				{
					Type:   batchv1.JobFailureTarget,
					Status: corev1.ConditionTrue,
				},
			},
			finished:          false,
			wantConditionType: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			finished, conditionType := JobFinished(&batchv1.Job{
				Status: batchv1.JobStatus{
					Conditions: tc.conditions,
				},
			})
			if diff := cmp.Diff(tc.finished, finished); diff != "" {
				t.Errorf("unexpected finished value (+got/-want): %s", diff)
			}
			if diff := cmp.Diff(tc.wantConditionType, conditionType); diff != "" {
				t.Errorf("unexpected condition type (+got/-want): %s", diff)
			}
		})
	}
}

func TestConstructJobsFromTemplate(t *testing.T) {
	var (
		jobSetName        = "test-jobset"
		replicatedJobName = "replicated-job"
		jobName           = "test-job"
		ns                = "default"
		topologyDomain    = "test-topology-domain"
	)

	tests := []struct {
		name      string
		js        *jobset.JobSet
		ownedJobs *childJobs
		want      []*batchv1.Job
	}{
		{
			name: "no jobs created",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(1).
					Obj()).Obj(),
			ownedJobs: &childJobs{
				active: []*batchv1.Job{
					testutils.MakeJob("test-jobset-replicated-job-0", ns).Obj(),
				},
			},
		},
		{
			name: "all jobs created",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(2).
					Obj()).Obj(),
			ownedJobs: &childJobs{},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-0",
					ns:                ns,
					replicas:          2,
					jobIdx:            0}).
					Suspend(false).Obj(),
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-1",
					ns:                ns,
					replicas:          2,
					jobIdx:            1}).
					Suspend(false).Obj(),
			},
		},
		{
			name: "one job created, one job not created (already active)",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(2).
					Obj()).Obj(),
			ownedJobs: &childJobs{
				active: []*batchv1.Job{
					testutils.MakeJob("test-jobset-replicated-job-0", ns).Obj(),
				},
			},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-1",
					ns:                ns,
					replicas:          2,
					jobIdx:            1}).
					Suspend(false).Obj(),
			},
		},
		{
			name: "one job created, one job not created (already succeeded)",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(2).
					Obj()).Obj(),
			ownedJobs: &childJobs{
				successful: []*batchv1.Job{
					testutils.MakeJob("test-jobset-replicated-job-0", ns).Obj(),
				},
			},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-1",
					ns:                ns,
					replicas:          2,
					jobIdx:            1}).
					Suspend(false).Obj(),
			},
		},
		{
			name: "one job created, one job not created (already failed)",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(2).
					Obj()).Obj(),
			ownedJobs: &childJobs{
				failed: []*batchv1.Job{
					testutils.MakeJob("test-jobset-replicated-job-0", ns).Obj(),
				},
			},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-1",
					ns:                ns,
					replicas:          2,
					jobIdx:            1}).
					Suspend(false).Obj(),
			},
		},
		{
			name: "one job created, one job not created (marked for deletion)",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(2).
					Obj()).Obj(),
			ownedJobs: &childJobs{
				delete: []*batchv1.Job{
					testutils.MakeJob("test-jobset-replicated-job-0", ns).Obj(),
				},
			},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-1",
					ns:                ns,
					replicas:          2,
					jobIdx:            1}).
					Suspend(false).Obj(),
			},
		},
		{
			name: "multiple replicated jobs",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-A").
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(1).
					Obj()).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-B").
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(2).
					Obj()).
				Obj(),
			ownedJobs: &childJobs{
				active: []*batchv1.Job{
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-B",
						jobName:           "test-jobset-replicated-job-B-0",
						ns:                ns,
						replicas:          2,
						jobIdx:            0}).
						Suspend(false).Obj(),
				},
			},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: "replicated-job-A",
					jobName:           "test-jobset-replicated-job-A-0",
					ns:                ns,
					replicas:          1,
					jobIdx:            0}).
					Suspend(false).Obj(),
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: "replicated-job-B",
					jobName:           "test-jobset-replicated-job-B-1",
					ns:                ns,
					replicas:          2,
					jobIdx:            1}).
					Suspend(false).Obj(),
			},
		},
		{
			name: "exclusive affinities",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).
						SetAnnotations(map[string]string{jobset.ExclusiveKey: topologyDomain}).
						Obj()).
					Replicas(1).
					Obj()).
				Obj(),
			ownedJobs: &childJobs{},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-0",
					ns:                ns,
					replicas:          1,
					jobIdx:            0,
					topology:          topologyDomain}).
					Suspend(false).Obj(),
			},
		},
		{
			name: "pod dns hostnames enabled",
			js: testutils.MakeJobSet(jobSetName, ns).
				EnableDNSHostnames(true).
				NetworkSubdomain(jobSetName).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Subdomain(jobSetName).
					Replicas(1).
					Obj()).
				Obj(),
			ownedJobs: &childJobs{},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-0",
					ns:                ns,
					replicas:          1,
					jobIdx:            0}).
					Suspend(false).
					Subdomain(jobSetName).Obj(),
			},
		},
		{
			name: "suspend job set",
			js: testutils.MakeJobSet(jobSetName, ns).
				Suspend(true).
				EnableDNSHostnames(true).
				NetworkSubdomain(jobSetName).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Subdomain(jobSetName).
					Replicas(1).
					Obj()).
				Obj(),
			ownedJobs: &childJobs{},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-0",
					ns:                ns,
					replicas:          1,
					jobIdx:            0}).
					Suspend(true).
					Subdomain(jobSetName).Obj(),
			},
		},
		{
			name: "resume job set",
			js: testutils.MakeJobSet(jobSetName, ns).
				Suspend(false).
				EnableDNSHostnames(true).
				NetworkSubdomain(jobSetName).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Subdomain(jobSetName).
					Replicas(1).
					Obj()).
				Obj(),
			ownedJobs: &childJobs{},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:        jobSetName,
					replicatedJobName: replicatedJobName,
					jobName:           "test-jobset-replicated-job-0",
					ns:                ns,
					replicas:          1,
					jobIdx:            0}).
					Suspend(false).
					Subdomain(jobSetName).Obj(),
			},
		},
		{
			name: "node selector exclusive placement strategy enabled",
			js: testutils.MakeJobSet(jobSetName, ns).
				EnableDNSHostnames(true).
				NetworkSubdomain(jobSetName).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).
						SetAnnotations(map[string]string{
							jobset.ExclusiveKey:            topologyDomain,
							jobset.NodeSelectorStrategyKey: "true",
						}).
						Obj()).
					Subdomain(jobSetName).
					Replicas(1).
					Obj()).
				Obj(),
			ownedJobs: &childJobs{},
			want: []*batchv1.Job{
				makeJob(&makeJobArgs{
					jobSetName:           jobSetName,
					replicatedJobName:    replicatedJobName,
					jobName:              "test-jobset-replicated-job-0",
					ns:                   ns,
					replicas:             1,
					jobIdx:               0,
					topology:             topologyDomain,
					nodeSelectorStrategy: true}).
					Suspend(false).
					Subdomain(jobSetName).
					NodeSelector(map[string]string{
						jobset.NamespacedJobKey: namespacedJobName(ns, "test-jobset-replicated-job-0"),
					}).
					Tolerations([]corev1.Toleration{
						{
							Key:      jobset.NoScheduleTaintKey,
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
					}).Obj(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []*batchv1.Job
			for _, rjob := range tc.js.Spec.ReplicatedJobs {
				jobs, err := constructJobsFromTemplate(tc.js, &rjob, tc.ownedJobs)
				if err != nil {
					t.Errorf("constructJobsFromTemplate() error = %v", err)
					return
				}
				got = append(got, jobs...)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("constructJobsFromTemplate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateConditions(t *testing.T) {
	var (
		jobSetName        = "test-jobset"
		replicatedJobName = "replicated-job"
		jobName           = "test-job"
		ns                = "default"
	)

	tests := []struct {
		name           string
		js             *jobset.JobSet
		conditions     []metav1.Condition
		newCondition   metav1.Condition
		expectedUpdate bool
	}{
		{
			name: "no condition",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(1).
					Obj()).Obj(),
			newCondition:   metav1.Condition{},
			conditions:     []metav1.Condition{},
			expectedUpdate: false,
		},
		{
			name: "do not update if false",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(1).
					Obj()).Obj(),
			newCondition:   metav1.Condition{Status: metav1.ConditionFalse, Type: string(jobset.JobSetSuspended), Reason: "JobsResumed"},
			conditions:     []metav1.Condition{},
			expectedUpdate: false,
		},
		{
			name: "update if condition is true",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(1).
					Obj()).Obj(),
			newCondition:   metav1.Condition{Status: metav1.ConditionTrue, Type: string(jobset.JobSetSuspended), Reason: "JobsResumed"},
			conditions:     []metav1.Condition{},
			expectedUpdate: true,
		},

		{
			name: "suspended",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(1).
					Obj()).Obj(),
			newCondition:   metav1.Condition{Status: metav1.ConditionTrue, Type: string(jobset.JobSetSuspended), Reason: "JobsSuspended"},
			conditions:     []metav1.Condition{},
			expectedUpdate: true,
		},
		{
			name: "resumed",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(1).
					Obj()).Obj(),
			newCondition:   metav1.Condition{Type: string(jobset.JobSetSuspended), Reason: "JobsResumed", Status: metav1.ConditionStatus(corev1.ConditionFalse)},
			conditions:     []metav1.Condition{{Type: string(jobset.JobSetSuspended), Reason: "JobsSuspended", Status: metav1.ConditionStatus(corev1.ConditionTrue)}},
			expectedUpdate: true,
		},
		{
			name: "duplicateComplete",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob(replicatedJobName).
					Job(testutils.MakeJobTemplate(jobName, ns).Obj()).
					Replicas(1).
					Obj()).Obj(),
			newCondition:   metav1.Condition{Type: string(jobset.JobSetCompleted), Message: "Jobs completed", Reason: "JobsCompleted", Status: metav1.ConditionTrue},
			conditions:     []metav1.Condition{{Type: string(jobset.JobSetCompleted), Message: "Jobs completed", Reason: "JobsCompleted", Status: metav1.ConditionTrue}},
			expectedUpdate: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsWithConditions := tc.js
			jsWithConditions.Status.Conditions = tc.conditions
			gotUpdate := updateCondition(jsWithConditions, tc.newCondition)
			if gotUpdate != tc.expectedUpdate {
				t.Errorf("updateCondition return mismatch")
			}
		})
	}
}

func TestCalculateReplicatedJobStatuses(t *testing.T) {
	var (
		jobSetName = "test-jobset"
		ns         = "default"
	)
	tests := []struct {
		name     string
		js       *jobset.JobSet
		jobs     childJobs
		expected []jobset.ReplicatedJobStatus
	}{
		{
			name: "partial jobs are ready, no succeeded jobs",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-1").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(1).
					Obj()).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-2").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(3).
					Obj()).Obj(),
			jobs: childJobs{
				active: []*batchv1.Job{
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-1",
						jobName:           "test-jobset-replicated-job-1-test-job-0",
						ns:                ns,
						replicas:          1,
						jobIdx:            0}).
						Parallelism(1).
						Completions(2).
						Ready(1).
						Succeeded(1).Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-0",
						ns:                ns,
						replicas:          3,
						jobIdx:            0}).
						Parallelism(5).
						Ready(2).
						Succeeded(3).Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-1",
						ns:                ns,
						replicas:          3,
						jobIdx:            0}).
						Parallelism(3).
						Completions(2).
						Ready(1).
						Succeeded(1).Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-2",
						ns:                ns,
						replicas:          3,
						jobIdx:            0}).
						Parallelism(2).
						Completions(3).
						Ready(2).
						Succeeded(1).Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-3",
						ns:                ns,
						replicas:          3,
						jobIdx:            0}).
						Parallelism(4).
						Completions(5).
						Ready(2).
						Succeeded(1).Obj(),
				},
			},
			expected: []jobset.ReplicatedJobStatus{
				{
					Name:      "replicated-job-1",
					Ready:     1,
					Succeeded: 0,
				},
				{
					Name:      "replicated-job-2",
					Ready:     3,
					Succeeded: 0,
				},
			},
		},
		{
			name: "no jobs created",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-1").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(1).
					Obj()).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-2").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(3).
					Obj()).Obj(),
			expected: []jobset.ReplicatedJobStatus{
				{
					Name:      "replicated-job-1",
					Ready:     0,
					Succeeded: 0,
				},
				{
					Name:      "replicated-job-2",
					Ready:     0,
					Succeeded: 0,
				},
			},
		},
		{
			name: "partial jobs created",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-1").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(1).
					Obj()).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-2").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(3).
					Obj()).Obj(),
			jobs: childJobs{
				active: []*batchv1.Job{
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-0",
						ns:                ns,
						replicas:          3,
						jobIdx:            0}).
						Parallelism(5).
						Ready(2).
						Succeeded(3).Obj(),
				},
			},
			expected: []jobset.ReplicatedJobStatus{
				{
					Name:      "replicated-job-1",
					Ready:     0,
					Succeeded: 0,
				},
				{
					Name:      "replicated-job-2",
					Ready:     1,
					Succeeded: 0,
				},
			},
		},
		{
			name: "no ready jobs, only succeeded jobs",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-1").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(1).
					Obj()).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-2").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(3).
					Obj()).Obj(),
			jobs: childJobs{
				successful: []*batchv1.Job{
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-0"}).Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-1",
						jobName:           "test-jobset-replicated-job-1-test-job-0"}).Obj(),
				},
			},
			expected: []jobset.ReplicatedJobStatus{
				{
					Name:      "replicated-job-1",
					Ready:     0,
					Succeeded: 1,
				},
				{
					Name:      "replicated-job-2",
					Ready:     0,
					Succeeded: 1,
				},
			},
		},
		{
			name: "no ready jobs, only failed jobs",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-1").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(1).
					Obj()).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-2").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(3).
					Obj()).Obj(),
			jobs: childJobs{
				failed: []*batchv1.Job{
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-0"}).Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-1",
						jobName:           "test-jobset-replicated-job-1-test-job-0"}).Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-1",
						jobName:           "test-jobset-replicated-job-1-test-job-1"}).Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-1",
						jobName:           "test-jobset-replicated-job-1-test-job-2"}).Obj(),
				},
			},
			expected: []jobset.ReplicatedJobStatus{
				{
					Name:   "replicated-job-1",
					Ready:  0,
					Failed: 3,
				},
				{
					Name:   "replicated-job-2",
					Ready:  0,
					Failed: 1,
				},
			},
		},
		{
			name: "active jobs",
			js: testutils.MakeJobSet(jobSetName, ns).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-1").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(1).
					Obj()).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-2").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(3).
					Obj()).Obj(),
			jobs: childJobs{
				active: []*batchv1.Job{
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-1",
						jobName:           "test-jobset-replicated-job-2-test-job-0"}).
						Parallelism(5).
						Active(1).
						Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-0"}).
						Parallelism(5).
						Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-1"}).
						Parallelism(1).
						Active(1).
						Obj(),
				},
			},
			expected: []jobset.ReplicatedJobStatus{
				{
					Name:   "replicated-job-1",
					Ready:  0,
					Active: 1,
				},
				{
					Name:   "replicated-job-2",
					Ready:  0,
					Active: 1,
				},
			},
		},
		{
			name: "suspended jobs",
			js: testutils.MakeJobSet(jobSetName, ns).Suspend(true).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-1").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(1).
					Obj()).
				ReplicatedJob(testutils.MakeReplicatedJob("replicated-job-2").
					Job(testutils.MakeJobTemplate("test-job", ns).Obj()).
					Replicas(3).
					Obj()).Obj(),
			jobs: childJobs{
				active: []*batchv1.Job{
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-1",
						jobName:           "test-jobset-replicated-job-2-test-job-0"}).
						Parallelism(5).
						Suspend(true).
						Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-0"}).
						Parallelism(5).
						Obj(),
					makeJob(&makeJobArgs{
						jobSetName:        jobSetName,
						replicatedJobName: "replicated-job-2",
						jobName:           "test-jobset-replicated-job-2-test-job-1"}).
						Parallelism(1).
						Suspend(true).
						Obj(),
				},
			},
			expected: []jobset.ReplicatedJobStatus{
				{
					Name:      "replicated-job-1",
					Ready:     0,
					Suspended: 1,
				},
				{
					Name:      "replicated-job-2",
					Ready:     0,
					Suspended: 1,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := JobSetReconciler{Client: (fake.NewClientBuilder()).Build()}
			statuses := r.calculateReplicatedJobStatuses(context.TODO(), tc.js, &tc.jobs)
			var less interface{} = func(a, b jobset.ReplicatedJobStatus) bool {
				return a.Name < b.Name
			}
			if diff := cmp.Diff(tc.expected, statuses, cmpopts.SortSlices(less)); diff != "" {
				t.Errorf("calculateReplicatedJobStatuses() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type makeJobArgs struct {
	jobSetName           string
	replicatedJobName    string
	jobName              string
	ns                   string
	replicas             int
	jobIdx               int
	restarts             int
	topology             string
	nodeSelectorStrategy bool
}

func makeJob(args *makeJobArgs) *testutils.JobWrapper {
	labels := map[string]string{
		jobset.JobSetNameKey:         args.jobSetName,
		jobset.ReplicatedJobNameKey:  args.replicatedJobName,
		jobset.ReplicatedJobReplicas: strconv.Itoa(args.replicas),
		jobset.JobIndexKey:           strconv.Itoa(args.jobIdx),
		RestartsKey:                  strconv.Itoa(args.restarts),
		jobset.JobKey:                jobHashKey(args.ns, args.jobName),
	}
	annotations := map[string]string{
		jobset.JobSetNameKey:         args.jobSetName,
		jobset.ReplicatedJobNameKey:  args.replicatedJobName,
		jobset.ReplicatedJobReplicas: strconv.Itoa(args.replicas),
		jobset.JobIndexKey:           strconv.Itoa(args.jobIdx),
		RestartsKey:                  strconv.Itoa(args.restarts),
		jobset.JobKey:                jobHashKey(args.ns, args.jobName),
	}
	// Only set exclusive key if we are using exclusive placement per topology.
	if args.topology != "" {
		annotations[jobset.ExclusiveKey] = args.topology
		// Exclusive placement topology domain must be set in order to use the node selector strategy.
		if args.nodeSelectorStrategy {
			annotations[jobset.NodeSelectorStrategyKey] = "true"
		}
	}
	jobWrapper := testutils.MakeJob(args.jobName, args.ns).
		JobLabels(labels).
		JobAnnotations(annotations).
		PodLabels(labels).
		PodAnnotations(annotations)
	return jobWrapper
}
