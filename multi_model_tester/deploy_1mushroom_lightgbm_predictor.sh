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
cat <<EOF |kubectl apply -f - >>./deploy.log
apiVersion: serving.kserve.io/v1alpha1
kind: Predictor
metadata:
  name: ${name}
spec:
  modelType:
    name: lightgbm
  path: lightgbm
  storage:
    s3:
      secretKey: localMinIO
EOF
  status=""
  try=1
  while [ "$status" != "Loaded" ]
  do
    if [ $try -eq $RETRY ]
    then
      break
    fi
    status=`kubectl get predictor ${name} -o jsonpath='{.status.activeModelState}'`
    sleep 1
    try=`expr $try + 1`
  done

  stopT=`date +%H:%M:%S`
  stopEpoc=`date +%s`
  elps=`expr $stopEpoc - $startEpoc`
  echo name,$name,status,$status,start,$startT,end,$stopT,elps,$elps 2>&1
done
