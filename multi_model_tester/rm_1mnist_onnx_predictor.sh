#!/bin/bash

# if [ $# -ne 1 ]
# then
#   echo "USAGE: $0 [name]"
#   echo "     name - name of the predictor"
#   exit
# fi

list=$@
RETRY=20
for name in $list
do
  startT=`date +%H:%M:%S`
  startEpoc=`date +%s`
cat <<EOF |kubectl delete -f - >>./deploy.log
apiVersion: serving.kserve.io/v1alpha1
kind: Predictor
metadata:
  name: ${name}
spec:
  modelType:
    name: onnx
  path: onnx/mnist.onnx
  storage:
    s3:
      secretKey: localMinIO
EOF
  sleep 1
  stopT=`date +%H:%M:%S`
  stopEpoc=`date +%s`
  elps=`expr $stopEpoc - $startEpoc`
  echo op,rm,name,$name,status,$status,start,$startT,end,$stopT,elps,$elps 2>&1
done
