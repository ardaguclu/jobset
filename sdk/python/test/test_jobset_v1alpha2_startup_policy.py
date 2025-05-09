# coding: utf-8

"""
    JobSet SDK

    Python SDK for the JobSet API

    The version of the OpenAPI document: v0.1.4
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


import unittest

from jobset.models.jobset_v1alpha2_startup_policy import JobsetV1alpha2StartupPolicy

class TestJobsetV1alpha2StartupPolicy(unittest.TestCase):
    """JobsetV1alpha2StartupPolicy unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional) -> JobsetV1alpha2StartupPolicy:
        """Test JobsetV1alpha2StartupPolicy
            include_optional is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # uncomment below to create an instance of `JobsetV1alpha2StartupPolicy`
        """
        model = JobsetV1alpha2StartupPolicy()
        if include_optional:
            return JobsetV1alpha2StartupPolicy(
                startup_policy_order = ''
            )
        else:
            return JobsetV1alpha2StartupPolicy(
                startup_policy_order = '',
        )
        """

    def testJobsetV1alpha2StartupPolicy(self):
        """Test JobsetV1alpha2StartupPolicy"""
        # inst_req_only = self.make_instance(include_optional=False)
        # inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()
