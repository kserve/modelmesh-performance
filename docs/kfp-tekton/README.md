## Setup KubeFlow Pipelines with Tekton

### Install Standalone Kubeflow Pipelines with Tekton
To install a minimum standalone setup for Kubeflow Pipelines with Tekton, visit the [Standalone Kubeflow Pipelines with Tekton Backend Deployment](https://github.com/kubeflow/kfp-tekton/blob/master/guides/kfp_tekton_install.md#standalone-kubeflow-pipelines-with-tekton-backend-deployment)

Once the standalone Kubeflow Pipelines with Tekton is installed, visit the [KFP-Tekton](https://github.com/kubeflow/kfp-tekton/tree/master/guides/kfp-user-guide) user guides to learn about how to upload and run a simple pipeline.

### Run Pipelines To Trigger KServe Performance Tests
After installing the standalone setup for Kubeflow Pipelines with Tekton from above, upload the pipelines below and start deploying and testing models with the latest KServe and ModelMesh deployments.

- [Deployment and payload demo on the simple string TF model](/multi_model_tester/modelmesh.yaml)
- [Send payloads to simple string TF, ONNX Minst, and Pytorch Cifar10 models](/multi_model_tester/modelmesh-payload.yaml)
- [Model removal without logs](/multi_model_tester/modelmesh-delete.yaml)
