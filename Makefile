# Copyright 2022 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Run unit tests
test:
	python3 -m unittest unit_test/k6/test_renderer.py

# Build howitzer docker images
build.howitzer:
	docker build -t k6_test:dev -f perf_test/docker/Dockerfile.slim .

# Deploy the resources used with the performance tests
deploy.predictor:
#	kubectl apply -f ./perf_test/k8s/example-mnist-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-example-mnist-predictor.yaml

#	oc apply -f ./perf_test/k8s/wml-serving-snapmlmnist-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-keras-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-pytorch-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-lightgbm-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-onnx-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-tensorflow-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-mleap-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-pmml-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-spark-predictor.yaml
#	oc apply -f ./perf_test/k8s/wml-serving-xgboost-predictor.yaml
	oc apply -f ./perf_test/k8s/wml-serving-sklearnMLserver-predictor.yaml

# Run perf-test from the local machine
run.howitzer-local:
#	CONFIG_FILE=./perf_test/configs/inference-test-lightgbm.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-onnx.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-tensorflow.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-pytorch.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test_skmnist.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-xgboost.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-mleap.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-spark.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-pmml.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-keras.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
	CONFIG_FILE=./perf_test/configs/inference-test-sklearn.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh
#	CONFIG_FILE=./perf_test/configs/inference-test-snapml.json TEMPLATE_DIR=./k6_test/ EXIT_ON_FAILURE=true ./perf_test/runHowitzer.sh

# snapml exists? sklearn == snapml? how to do mleap and pmml


# Run perf-test inside a k8s cluster using pre-built image
run.perf-test: run.delete-perf-test-job deploy.predictor
	kubectl apply -f ./perf_test/k8s/howitzer_k6_test-k8s.yaml
	# kubectl wait po -l='job-name=perf-test-job' --for=condition=READY=true --timeout=300s
	# kubectl logs -f `kubectl get po -l='job-name=perf-test-job' -o jsonpath='{.items[*].metadata.name}'`
run.delete-perf-test-job:
	kubectl delete -f ./perf_test/k8s/howitzer_k6_test-k8s.yaml
