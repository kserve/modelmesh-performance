#!/bin/bash

list=$@
RETRY=20
for name in $list
do
  startT=`date +%H:%M:%S`
  startEpoc=`date +%s`
cat <<EOF |kubectl delete -f - >>./rm.log
apiVersion: serving.kserve.io/v1alpha1
kind: Predictor
metadata:
  name: ${name}
spec:
  modelType:
    name: tensorflow
  path: tensorflow/simple_string
  storage:
    s3:
      secretKey: localMinIO
EOF
  status=$?
  sleep 1
  stopT=`date +%H:%M:%S`
  stopEpoc=`date +%s`
  elps=`expr $stopEpoc - $startEpoc`
  echo name,$name,status,$status,start,$startT,end,$stopT,elps,$elps 2>&1
done
