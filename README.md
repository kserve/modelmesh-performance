# ModelMesh Metrics
ModelMesh Performance Scripts, Dashboard and Pipelines.
Follow the instruction to run gRPC inference requests benchmark on a ModelMesh Serving instance on a Kubernetes cluster:
1. [Setup ModelMesh Serving](#Setup-ModelMesh-Serving)
2. [Create Example Models](#Model-Deployment)
3. [Run gRPC Inference Benchmark](#Run-Inference-Requests)
4. [Monitoring](#Monitoring)
5. [Automation with KubeFlow Pipeline](#Test-with-KubeFlow-Pipeline)
## Setup ModelMesh Serving
To quickly stand up a ModelMesh Serving instance on a Kubernetes Cluster, please see [Quickstart](https://github.com/kserve/modelmesh-serving/blob/main/docs/quickstart.md) for detail.

## Model Deployment
Switch to the ModelMesah Serving namespace, assuming `modelmesh-serving` :
```bash
kubectl config set-context --current --namespace modelmesh-serving
```
We have included [scripts](./docs/README.md#Using-deployment-scripts) to deploy models concurrently to a ModelMesh Serving instance.
Deploy 10 simple-string tensorflow models at concurrency of 5:
```bash
cd multi_model_tester/
./deployNpredictors.sh 5 simple-string-tf 1 10 deploy_1simple_string_tf_predictor.sh
```
Verify the 10 simple-string-tf-* models are `Loaded`:
```bash
kubectl get predictors |grep simple-string-tf
```

## Run Inference Requests
On a local machine, assuming a separate terminal is running `kubectl port-forward svc/modelmesh-serving 8033`.
Run the [multi_model_tester](./docs/README.md#build-the-tester-binary) to send inference requests to the 10 simple-string-tf-* models.
```bash
./multi_model_test -ma "SimpleStringTF" -npm 10 -qps 100 -dur 10
```

Running the driver from inside the Kubernetes cluster and namespace where the Modelmesh Serving instance is installed:
```bash
KUBECTL_COMMAND_HEADERS=false kubectl run -it --rm modelmesh-payload --restart=Never --image=aipipeline/modelmesh-payload:latest --command --namespace modelmesh-serving -- ./multi_model_test -ma 'SimpleStringTF' -u dns:///modelmesh-serving.modelmesh-serving:8033 -npm '10'  -qps '100' -dur '10'
```

Depending on where the driver is run, output should be similar to:
```
QPS: 100
dur: 9.993041819
earliest: 76448
latest: 9990014958
end: 9992985483
reqs: 1000
success: 1.000000
p50: 3.58451ms
p95: 10.701875ms
p99: 29.714907ms
mean: 4.733809ms
max: 37.563672ms
```

## Monitoring
Monitoring ModelMesh Serving metrics using Prometheus and Grafana Dashboard is highly recommended. See [monitoring](./docs/monitoring/README.md##Setup-Prometheus-Operator) for detail.


## Test with KubeFlow Pipeline
Deploying models and sending inference request can be automated using KubeFlow Pipeline. See [Setup KubeFlow Tekton](./docs/kfp-tekton/README.md##Setup-KubeFlow-Tekton) for detail.



