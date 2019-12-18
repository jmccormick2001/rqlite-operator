#!/bin/bash

VOLNAME="/mylpv"
DIRNAME="vol1"
NS=rqnamespace

sudo mkdir -p $VOLNAME/$DIRNAME 
sudo chcon -Rt svirt_sandbox_file_t $VOLNAME/$DIRNAME
sudo chmod 777 $VOLNAME/$DIRNAME

kubectl create -f storageClass.yaml -n $NS
kubectl create -f persistentVolume.yaml -n $NS
kubectl create -f persistentVolumeClaim.yaml -n $NS
kubectl create -f http-pod.yaml -n $NS
cat > $VOLNAME/$DIRNAME/index.html << EOF
<html>
hello from the rqlite operator persistence test
</html>
EOF
