#!/bin/bash

VOLNAME="/mylpv"
DIRNAME="vol1"
NS=rqnamespace

kubectl delete -f storageClass.yaml -n $NS
kubectl delete -f persistentVolume.yaml -n $NS
kubectl delete -f persistentVolumeClaim.yaml -n $NS
kubectl delete -f http-pod.yaml -n $NS
