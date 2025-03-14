# coding: utf-8

"""
    JobSet SDK

    Python SDK for the JobSet API

    The version of the OpenAPI document: v0.1.4
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


import unittest

from jobset.models.io_k8s_api_batch_v1_pod_failure_policy import IoK8sApiBatchV1PodFailurePolicy

class TestIoK8sApiBatchV1PodFailurePolicy(unittest.TestCase):
    """IoK8sApiBatchV1PodFailurePolicy unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional) -> IoK8sApiBatchV1PodFailurePolicy:
        """Test IoK8sApiBatchV1PodFailurePolicy
            include_optional is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # uncomment below to create an instance of `IoK8sApiBatchV1PodFailurePolicy`
        """
        model = IoK8sApiBatchV1PodFailurePolicy()
        if include_optional:
            return IoK8sApiBatchV1PodFailurePolicy(
                rules = [
                    jobset.models.io/k8s/api/batch/v1/pod_failure_policy_rule.io.k8s.api.batch.v1.PodFailurePolicyRule(
                        action = '', 
                        on_exit_codes = jobset.models.io/k8s/api/batch/v1/pod_failure_policy_on_exit_codes_requirement.io.k8s.api.batch.v1.PodFailurePolicyOnExitCodesRequirement(
                            container_name = '', 
                            operator = '', 
                            values = [
                                56
                                ], ), 
                        on_pod_conditions = [
                            jobset.models.io/k8s/api/batch/v1/pod_failure_policy_on_pod_conditions_pattern.io.k8s.api.batch.v1.PodFailurePolicyOnPodConditionsPattern(
                                status = '', 
                                type = '', )
                            ], )
                    ]
            )
        else:
            return IoK8sApiBatchV1PodFailurePolicy(
                rules = [
                    jobset.models.io/k8s/api/batch/v1/pod_failure_policy_rule.io.k8s.api.batch.v1.PodFailurePolicyRule(
                        action = '', 
                        on_exit_codes = jobset.models.io/k8s/api/batch/v1/pod_failure_policy_on_exit_codes_requirement.io.k8s.api.batch.v1.PodFailurePolicyOnExitCodesRequirement(
                            container_name = '', 
                            operator = '', 
                            values = [
                                56
                                ], ), 
                        on_pod_conditions = [
                            jobset.models.io/k8s/api/batch/v1/pod_failure_policy_on_pod_conditions_pattern.io.k8s.api.batch.v1.PodFailurePolicyOnPodConditionsPattern(
                                status = '', 
                                type = '', )
                            ], )
                    ],
        )
        """

    def testIoK8sApiBatchV1PodFailurePolicy(self):
        """Test IoK8sApiBatchV1PodFailurePolicy"""
        # inst_req_only = self.make_instance(include_optional=False)
        # inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()
