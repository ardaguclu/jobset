# coding: utf-8

"""
    JobSet SDK

    Python SDK for the JobSet API

    The version of the OpenAPI document: v0.1.4
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


import unittest

from jobset.models.io_k8s_api_core_v1_container_resize_policy import IoK8sApiCoreV1ContainerResizePolicy

class TestIoK8sApiCoreV1ContainerResizePolicy(unittest.TestCase):
    """IoK8sApiCoreV1ContainerResizePolicy unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional) -> IoK8sApiCoreV1ContainerResizePolicy:
        """Test IoK8sApiCoreV1ContainerResizePolicy
            include_optional is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # uncomment below to create an instance of `IoK8sApiCoreV1ContainerResizePolicy`
        """
        model = IoK8sApiCoreV1ContainerResizePolicy()
        if include_optional:
            return IoK8sApiCoreV1ContainerResizePolicy(
                resource_name = '',
                restart_policy = ''
            )
        else:
            return IoK8sApiCoreV1ContainerResizePolicy(
                resource_name = '',
                restart_policy = '',
        )
        """

    def testIoK8sApiCoreV1ContainerResizePolicy(self):
        """Test IoK8sApiCoreV1ContainerResizePolicy"""
        # inst_req_only = self.make_instance(include_optional=False)
        # inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()
