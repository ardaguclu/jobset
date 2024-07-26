# coding: utf-8

"""
    JobSet SDK

    Python SDK for the JobSet API  # noqa: E501

    The version of the OpenAPI document: v0.1.4
    Generated by: https://openapi-generator.tech
"""


from __future__ import absolute_import

# Kubernetes imports
from kubernetes.client.models.v1_job_template_spec import V1JobTemplateSpec
import unittest
import datetime

import jobset
from jobset.models.jobset_v1alpha2_job_set_list import JobsetV1alpha2JobSetList  # noqa: E501
from jobset.rest import ApiException

class TestJobsetV1alpha2JobSetList(unittest.TestCase):
    """JobsetV1alpha2JobSetList unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional):
        """Test JobsetV1alpha2JobSetList
            include_option is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # model = jobset.models.jobset_v1alpha2_job_set_list.JobsetV1alpha2JobSetList()  # noqa: E501
        if include_optional :
            return JobsetV1alpha2JobSetList(
                api_version = '0', 
                items = [
                    jobset.models.jobset_v1alpha2_job_set.JobsetV1alpha2JobSet(
                        api_version = '0', 
                        kind = '0', 
                        metadata = None, 
                        spec = jobset.models.jobset_v1alpha2_job_set_spec.JobsetV1alpha2JobSetSpec(
                            coordinator = jobset.models.jobset_v1alpha2_coordinator.JobsetV1alpha2Coordinator(
                                job_index = 56, 
                                pod_index = 56, 
                                replicated_job = '0', ), 
                            failure_policy = jobset.models.jobset_v1alpha2_failure_policy.JobsetV1alpha2FailurePolicy(
                                max_restarts = 56, 
                                rules = [
                                    jobset.models.jobset_v1alpha2_failure_policy_rule.JobsetV1alpha2FailurePolicyRule(
                                        action = '0', 
                                        name = '0', 
                                        on_job_failure_reasons = [
                                            '0'
                                            ], 
                                        target_replicated_jobs = [
                                            '0'
                                            ], )
                                    ], ), 
                            managed_by = '0', 
                            network = jobset.models.jobset_v1alpha2_network.JobsetV1alpha2Network(
                                enable_dns_hostnames = True, 
                                publish_not_ready_addresses = True, 
                                subdomain = '0', ), 
                            replicated_jobs = [
                                jobset.models.jobset_v1alpha2_replicated_job.JobsetV1alpha2ReplicatedJob(
                                    name = '0', 
                                    replicas = 56, 
                                    template = V1JobTemplateSpec(), )
                                ], 
                            startup_policy = jobset.models.jobset_v1alpha2_startup_policy.JobsetV1alpha2StartupPolicy(
                                startup_policy_order = '0', ), 
                            success_policy = jobset.models.jobset_v1alpha2_success_policy.JobsetV1alpha2SuccessPolicy(
                                operator = '0', ), 
                            suspend = True, 
                            ttl_seconds_after_finished = 56, ), 
                        status = jobset.models.jobset_v1alpha2_job_set_status.JobsetV1alpha2JobSetStatus(
                            conditions = [
                                None
                                ], 
                            replicated_jobs_status = [
                                jobset.models.jobset_v1alpha2_replicated_job_status.JobsetV1alpha2ReplicatedJobStatus(
                                    active = 56, 
                                    failed = 56, 
                                    name = '0', 
                                    ready = 56, 
                                    succeeded = 56, 
                                    suspended = 56, )
                                ], 
                            restarts = 56, 
                            restarts_count_towards_max = 56, 
                            terminal_state = '0', ), )
                    ], 
                kind = '0', 
                metadata = None
            )
        else :
            return JobsetV1alpha2JobSetList(
                items = [
                    jobset.models.jobset_v1alpha2_job_set.JobsetV1alpha2JobSet(
                        api_version = '0', 
                        kind = '0', 
                        metadata = None, 
                        spec = jobset.models.jobset_v1alpha2_job_set_spec.JobsetV1alpha2JobSetSpec(
                            coordinator = jobset.models.jobset_v1alpha2_coordinator.JobsetV1alpha2Coordinator(
                                job_index = 56, 
                                pod_index = 56, 
                                replicated_job = '0', ), 
                            failure_policy = jobset.models.jobset_v1alpha2_failure_policy.JobsetV1alpha2FailurePolicy(
                                max_restarts = 56, 
                                rules = [
                                    jobset.models.jobset_v1alpha2_failure_policy_rule.JobsetV1alpha2FailurePolicyRule(
                                        action = '0', 
                                        name = '0', 
                                        on_job_failure_reasons = [
                                            '0'
                                            ], 
                                        target_replicated_jobs = [
                                            '0'
                                            ], )
                                    ], ), 
                            managed_by = '0', 
                            network = jobset.models.jobset_v1alpha2_network.JobsetV1alpha2Network(
                                enable_dns_hostnames = True, 
                                publish_not_ready_addresses = True, 
                                subdomain = '0', ), 
                            replicated_jobs = [
                                jobset.models.jobset_v1alpha2_replicated_job.JobsetV1alpha2ReplicatedJob(
                                    name = '0', 
                                    replicas = 56, 
                                    template = V1JobTemplateSpec(), )
                                ], 
                            startup_policy = jobset.models.jobset_v1alpha2_startup_policy.JobsetV1alpha2StartupPolicy(
                                startup_policy_order = '0', ), 
                            success_policy = jobset.models.jobset_v1alpha2_success_policy.JobsetV1alpha2SuccessPolicy(
                                operator = '0', ), 
                            suspend = True, 
                            ttl_seconds_after_finished = 56, ), 
                        status = jobset.models.jobset_v1alpha2_job_set_status.JobsetV1alpha2JobSetStatus(
                            conditions = [
                                None
                                ], 
                            replicated_jobs_status = [
                                jobset.models.jobset_v1alpha2_replicated_job_status.JobsetV1alpha2ReplicatedJobStatus(
                                    active = 56, 
                                    failed = 56, 
                                    name = '0', 
                                    ready = 56, 
                                    succeeded = 56, 
                                    suspended = 56, )
                                ], 
                            restarts = 56, 
                            restarts_count_towards_max = 56, 
                            terminal_state = '0', ), )
                    ],
        )

    def testJobsetV1alpha2JobSetList(self):
        """Test JobsetV1alpha2JobSetList"""
        inst_req_only = self.make_instance(include_optional=False)
        inst_req_and_optional = self.make_instance(include_optional=True)


if __name__ == '__main__':
    unittest.main()
