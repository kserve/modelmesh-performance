#!/bin/bash

list=$@
RETRY=20
for name in $list
do
  startT=`date +%H:%M:%S`
  startEpoc=`date +%s`
  kubectl delete predictor $name
  status=$?
  try=1
  while [ $status -ne 1 ]
  do
    if [ $try -eq $RETRY ]
    then
      break
    fi
    kubectl get predictor ${name} -o jsonpath='{.status.activeModelState}'
    status=$?

    sleep 1
    try=`expr $try + 1`
  done

  stopT=`date +%H:%M:%S`
  stopEpoc=`date +%s`
  elps=`expr $stopEpoc - $startEpoc`
  echo test,"rm",name,$name,status,$status,start,$startT,end,$stopT,elps,$elps 2>&1
done
