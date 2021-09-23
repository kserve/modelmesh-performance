#!/bin/bash

if [ $# -ne 4 ]
then
  echo "USAGE: $0 [parallel] [prefix] [startID] [endID]"
  echo "     parallel - number of predictors to be concurrelty to remove"
  echo "     prefix - prefix to use in the predictor name"
  echo "     startID - the starting ID of a prefix-ID"
  echo "     endID  - the ending ID of the predictor deploy"
  exit
fi


N=$1 #N is the number of predictors to deploy concurrently
startID=$3
endID=$4
PREFIX=$2

for i in `seq $startID $endID`
do
  list+=("$PREFIX-${i}")
done

echo startID is $startID
echo endID is $endID
echo parallel is $N
#echo ${list[*]}

#vm=(`nova list --all-tenants|grep $2|awk '{print $2}'`)
total=`expr $endID - $startID + 1`
if [ $total -eq 0 ]
then
  echo Nothing to create.
  exit 1
fi

factor=`expr $total / $N`
mod=`expr $total % $N`

#echo Total is $total
#echo Factor is $factor
#echo Modulo is $mod

#slice the list N times
for ((i=0; i < $N; i++))
do
  beg=`expr $i \* $factor` #beg index
  #echo beg is $beg
  slice=(${list[@]:$beg:$factor})
  #echo thread $i
  if [ $mod -gt 0 ]
  then
    tail=`expr $total - $mod`
    #echo The tail index is $tail
    #echo The tail content is ${list[$tail]}

    slice=("${slice[@]}" "${list[$tail]}")
    let "mod--"
  fi
  #echo ${slice[@]}
  ./rm_1simple_string_predictor.sh ${slice[@]}&
done
