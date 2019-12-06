rqliteimage:   
	sudo --preserve-env buildah bud -f ./rqlite-image/Dockerfile -t jemccorm/rqlite:v0.0.1 ./rqlite-image
	sudo --preserve-env buildah push --authfile /home/jeffmc/.docker/config.json jemccorm/rqlite:v0.0.1 docker://quay.io/jemccorm/rqlite:v0.0.1
	docker tag quay.io/jemccorm/rqlite:v0.0.1  jemccorm/rqlite:v0.0.1 
configmap:   
	kubectl delete configmap rq-config -n rqnamespace --ignore-not-found
	kubectl create configmap rq-config \
		--from-file=./templates/pod-template.json \
		--from-file=./templates/service-template.json \
		-n rqnamespace
testit:   
	kubectl create -f deploy/operator.yaml -n rqnamespace
	sleep 5
	kubectl create -f deploy/crds/rqcluster.example.com_v1alpha1_rqcluster_cr.yaml -n rqnamespace
testitlocal:   
	export OPERATOR_NAME=rqlite-operator
	operator-sdk up local --namespace=rqnamespace
cleanup:   
	kubectl delete rqclusters --all -n rqnamespace
	kubectl delete deploy --all -n rqnamespace
setup:
	kubectl create -n rqnamespace -f deploy/crds/rqcluster.example.com_rqclusters_crd.yaml
	kubectl create -n rqnamespace -f deploy/service_account.yaml
	kubectl create -n rqnamespace -f deploy/role.yaml
	kubectl create -n rqnamespace -f deploy/role_binding.yaml
operatorimage:   
	operator-sdk build quay.io/jemccorm/rqlite-operator:v0.0.1
pushoperatorimage:   
	docker push quay.io/jemccorm/rqlite-operator:v0.0.1
olmuninstall:   
	kubectl -n rqnamespace delete csv rqlite-operator.v0.0.1 --ignore-not-found
	kubectl -n rqnamespace delete operatorgroup rqlite-operator-group --ignore-not-found
	kubectl delete crd rqclusters.rqcluster.example.com --ignore-not-found
olminstall:   
	kubectl create -f deploy/olm-manual/operator-group.yaml -n rqnamespace
	kubectl create -f deploy/olm-catalog/rqlite-operator/0.0.1/rqlite-operator.v0.0.1.clusterserviceversion.yaml -n rqnamespace
	kubectl create -n rqnamespace -f deploy/crds/rqcluster.example.com_rqclusters_crd.yaml
	kubectl create -n rqnamespace -f deploy/service_account.yaml
	kubectl create -n rqnamespace -f deploy/role.yaml
	kubectl create -n rqnamespace -f deploy/role_binding.yaml

