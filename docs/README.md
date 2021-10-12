## Using deployment scripts
The scripts requires [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) and login your namespace where the ModelMesh Serving is installed, usually `modelmesh-serving` namespace. Below are examples of creating/removing models from the ModelMesh Serving instance using the included shell scripts:


```bash
cd multi_model_tester
```
```bash
# Deploy 5000 Mnist Onnx models named mnist-onnx-1 to mnist-onnx-5000 at concurrency of 10
./deployNpredictors.sh 10 mnist-onnx 1 5000 deploy_1mnist_onnx_predictor.sh
```
```bash
# Deploy 1000 simple string models named simple-string-tf-1000 to simple-string-2000 at concurrency of 30
./deployNpredictors.sh 30 simple-string-tf 1000 2000 deploy_1simple_string_tf_predictor.sh
```
To clean up the models:
```bash
# Remove Mnist Onnx models named mnist-onnx-1 to mnist-onnx-5000 at concurrency of 10
./rmNpredictors.sh 10 mnist-onnx 1 5000 rm_1mnist_onnx_predictor.sh
```
```bash
# Remove 1000 simple string models named simple-string-tf-1000 to simple-string-2000 at concurrency of 30
./rmNpredictors.sh 30 simple-string-tf 1000 2000 rm_1simple_string_tf_predictor.sh
```
## Using the multi_model_tester
The `multi_model_tester` is written in Go and utilizes [Trunks](https://github.com/straightdave/trunks.git) library that is modified to uses round_robbin loadbalancer policy and targets `dns:///modelmesh-serving:8033`.

The tester generates a mix of up to six heterogeneous model inference requests to stress the [inference request gRPC api](https://github.com/kserve/kserve/blob/master/docs/predict-api/v2/required_api.md#grpc) of ModelMesh Serving instance.

### Build the tester binary:
```bash
git clone https://github.com/tedhtchang/kserve-perf.git
cd kserve-perf/multi_model_tester
go build .
```

### Usage:
```bash
./multi_model_test -h
Usage of ./multi_model_test:
  -cp uint
    	Number of connections to create. Default: 1 (default 1)
  -debug
    	Sends 1 request and view response. Default: false
  -dur int
    	Test duration in seconds. Default: 1 (default 1)
  -ma string
    	List of different types models separate by space. Available options: SimpleStringTF, MnistSklearn, MushroomXgboost, CifarPytorch, MushroomLightgbm, MnistOnnx
  -npm int
    	Number of model name 1 to npm per model to generate (default 1)
  -qps int
    	Constant Queries Per Second to hold. Default: 1 (default 1)
  -u string
    	Inference Server URL. Default: localhost:8033 (default "localhost:8033")
  -wp uint
    	Number of worker pool. Default: 1 (default 1)
```
### Examples
Using the tester from inside same namespace to target the kube-dns address and port that is corresponding ModelMesh Serving service.

Target the Simple-String models named simple-string-tf-1 to simple-string-tf-2000 at 1000 queries for 60 seconds.
```bash
./multi_model_test -u dns:///modelmesh-serving:8033 -npm 2000 -ma "SimpleStringTF" -qps 1000 -dur 60
```
Target the Simple-String and MnistOnnx models named simple-string-tf-1 to simple-string-tf-2000 and mnist-onnx-1 to mnist-onnx-2000 at 1000 queries for 60 seconds.
```bash
./multi_model_test -u dns:///modelmesh-serving:8033 -npm 2000 -ma "SimpleStringTF MnistOnnx" -qps 1000 -dur 60
```
 Sometimes, it is useful to verify whether all tested the models can respond. In debug mode, the tester sends inference requests to Simple-String and MnistOnnx models named simple-string-tf-1 to simple-string-tf-2000 and mnist-onnx-1 to mnist-onnx-2000 sequentially and print the response.

```bash
./multi_model_test -u dns:///modelmesh-serving:8033 -npm 2000 -ma "SimpleStringTF MnistOnnx" -debug
```
