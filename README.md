# Introduction

# Environment Setup

## Setup ModelMesh Serving
To quickly stand up a ModelMesh Serving instance on a Kubernetes Cluster, please see [Quickstart](https://github.com/kserve/modelmesh-serving/blob/main/docs/quickstart.md) for detail.

## Model Deployment
Switch to the ModelMesah Serving namespace, assuming `modelmesh-serving` :
```bash
kubectl config set-context --current --namespace modelmesh-serving
```
We have included [scripts](./docs/Using-deployment-scripts) to deploy models concurrently to a ModelMesh Serving instance.
Deploy 10 simple-string tensorflow models at concurrency of 5:
```bash
cd multi_model_tester/
./deployNpredictors.sh 5 simple-string-tf 1 10 deploy_1simple_string_tf_predictor.sh
```
Verify you 10 simple-string-tf are `Loaded`:
```bash
kubectl get predictors |grep simple-string-tf
```

## Run Inference Requests
Using the [multi_model_tester](./docs/README.md##Using-the-multi_model_tester) to send inference requests to the 10 simple-string-tf-* models:
```bash
./multi_model_test -ma "SimpleStringTF" -npm 10 -qps 1 -dur 10
```
Depending on where the driver is run, the numbers may vary:
```
Created requests.
QPS: 1
dur: 9.097220044
earliest: 79102
latest: 9000005765
end: 9097189999
reqs: 10
success: 1.000000
p50: 107.056888ms
p95: 116.525302ms
p99: 116.525302ms
mean: 125.253537ms
max: 293.500996ms
```

## Monitoring (Optional)
Monitoring ModelMesh Serving metrics using Prometheus and Grafana Dashboard is highly recommended. See [monitoring](./docs/monitoring/README.md##Setup-Prometheus-Operator) for detail.


## Test with KubeFlow Pipeline (Optional)
Deploying models and sending inference request can be automated using KubeFlow Pipeline. See [Setup KubeFlow Tekton](./docs/kfp-tekton/README.md##Setup-KubeFlow-Tekton) for detail.
