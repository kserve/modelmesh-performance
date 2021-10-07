#!/bin/bash

if [ $# -ne 2 ]
then
  echo "USAGE: $0 [num_of_predictors] [model_type]"
  echo "     num_of_predictors - number of additional predictors to be created"
  echo "     model_type - one of the model types to be deployed as predictors"
  exit
fi


N=$1 #N is the number of predictors to deploy concurrently
model_name=$2

if [ "${model_name}" == "SimpleStringTF" ]
then
  script="deploy_1simple_string_tf_predictor.sh"
  PREFIX="simple-string-tf"
fi

if [ "${model_name}" == "MnistSklearn" ]
then
  script="deploy_1mnist_sklearn_predictor.sh"
  PREFIX="mnist-sklearn"
fi

if [ "${model_name}" == "MushroomXgboost" ]
then
  script="deploy_1mushroom_xgboost_predictor.sh"
  PREFIX="mushroom-xgboost"
fi

if [ "${model_name}" == "CifarPytorch" ]
then
  script="deploy_1cifar_pytorch_predictor.sh"
  PREFIX="cifar-pytorch"
fi

if [ "${model_name}" == "MushroomLightgbm" ]
then
  script="deploy_1mushroom_lightgbm_predictor.sh"
  PREFIX="mushroom-lightgbm"
fi

if [ "${model_name}" == "MnistOnnx" ]
then
  script="deploy_1mnist_onnx_predictor.sh"
  PREFIX="mnist-onnx"
fi

if [ ! -f "${script}" ]
then
  echo "$model_name model does not exist."
  exit 1
fi

model_num=$(kubectl get predictor -n modelmesh-serving | grep ${PREFIX} |  awk '{print $1}' | sed s/${PREFIX}-//g | sort -n -r | awk '{print $NF; exit}')

if [ -z "$model_num" ]
then
  model_num=0
fi

startID=$(( 1 + $model_num ))
endID=$(( $N + $model_num ))

./deployNpredictors.sh 10 $PREFIX $startID $endID $script
