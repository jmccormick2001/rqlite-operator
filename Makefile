NS = rqnamespace
rqliteimage:   
	sudo --preserve-env buildah bud -f ./rqlite-image/Dockerfile -t quay.io/jemccorm/rqlite:v0.0.2 ./rqlite-image
	sudo --preserve-env buildah push --authfile /home/jeffmc/.docker/config.json jemccorm/rqlite:v0.0.2 docker://quay.io/jemccorm/rqlite:v0.0.2
	docker tag quay.io/jemccorm/rqlite:v0.0.2  jemccorm/rqlite:v0.0.2 
configmap:   
	kubectl delete configmap rq-config -n $(NS) --ignore-not-found
	kubectl create configmap rq-config \
		--from-file=./templates/pod-template.json \
		-n $(NS)
testit:   
	@echo $(NS) is the namespace
	kubectl create -f deploy/operator.yaml -n $(NS)
	kubectl create -f deploy/operator2.yaml -n $(NS)
	sleep 5
	kubectl create -f deploy/crds/rqcluster.example.com_v1alpha1_rqcluster_cr.yaml -n $(NS)
testitlocal:   
	export OPERATOR_NAME=rqlite-operator
	operator-sdk up local --namespace=$(NS)
cleanup:   
	@echo $(NS) is the namespace
	kubectl delete rqclusters --all -n $(NS)
	kubectl delete deploy --all -n $(NS)
setup:
	@echo $(NS) is the namespace
	kubectl create -n $(NS) -f deploy/crds/rqcluster.example.com_rqclusters_crd.yaml
	kubectl create -n $(NS) -f deploy/service_account.yaml
	kubectl create -n $(NS) -f deploy/role.yaml
	kubectl create -n $(NS) -f deploy/role_binding.yaml
operatorimage:   
	operator-sdk build quay.io/jemccorm/rqlite-operator:v0.0.2
pushoperatorimage:   
	docker push quay.io/jemccorm/rqlite-operator:v0.0.2
olmuninstall:   
	kubectl -n $(NS) delete csv rqlite-operator.v0.0.1 --ignore-not-found
	kubectl -n $(NS) delete operatorgroup rqlite-operator-group --ignore-not-found
	kubectl delete crd rqclusters.rqcluster.example.com --ignore-not-found
olminstall:   
	kubectl create -f deploy/olm-manual/operator-group.yaml -n $(NS)
	kubectl create -f deploy/olm-catalog/rqlite-operator/0.0.1/rqlite-operator.v0.0.1.clusterserviceversion.yaml -n $(NS)
	kubectl create -n $(NS) -f deploy/crds/rqcluster.example.com_rqclusters_crd.yaml
	kubectl create -n $(NS) -f deploy/service_account.yaml
	kubectl create -n $(NS) -f deploy/role.yaml
	kubectl create -n $(NS) -f deploy/role_binding.yaml

