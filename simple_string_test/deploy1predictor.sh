#!/bin/bash

if [ $# -ne 1 ]
then
  echo "USAGE: $0 [name] [prefix]"
  echo "     name - name of the predictor"
  exit
fi

name="example-mnist-predictor-$1"
startT=`date +%H:%M:%S`
cat <<EOF |kubectl apply -f - >>./deploy.log
apiVersion: serving.kserve.io/v1alpha1
kind: Predictor
metadata:
  name: ${name}
spec:
  modelType:
    name: sklearn
  path: sklearn/mnist-svm.joblib
  storage:
    s3:
      secretKey: localMinIO
EOF
status=""
# while [ "$status" != "Loaded" ]
# do
#   status=`kubectl get predictor ${name} -o jsonpath='{.status.activeModelState}'`
#   sleep 1
# done
stopT=`date +%H:%M:%S`
echo name,$name,start,$startT,end,$stopT
