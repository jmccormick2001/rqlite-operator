NS = rq
IMAGEUSER = someuser
rqo:   
	go build -o bin/rqo github.com/jmccormick2001/rqlite-operator/pkg/cli
rqliteimage:   
	sudo --preserve-env buildah bud -f ./rqlite-image/Dockerfile -t quay.io/$(IMAGEUSER)/rqlite:v0.0.2 ./rqlite-image
	sudo --preserve-env buildah push --authfile /home/jeffmc/.docker/config.json $(IMAGEUSER)/rqlite:v0.0.2 docker://quay.io/$(IMAGEUSER)/rqlite:v0.0.2
test:   
	@echo $(NS) is the namespace
	kubectl create -f deploy/operator.yaml -n $(NS)
	kubectl create -f deploy/operator2.yaml -n $(NS)
	sleep 5
#	kubectl create -n $(NS) -f deploy/crds/rqcluster.example.com_v1alpha1_rqcluster_cr.yaml
testitlocal:   
	export OPERATOR_NAME=rqlite-operator
	operator-sdk up local --namespace=$(NS)
clean:   
	@echo $(NS) is the namespace
	kubectl delete rqclusters --all -n $(NS)
	kubectl delete deploy --all -n $(NS)
	kubectl delete namespace $(NS)
	kubectl delete crd rqclusters.rqcluster.example.com
verify: 
	kubectl -n $(NS) get deploy
	kubectl -n $(NS) get pod
	kubectl -n $(NS) get svc
	kubectl -n $(NS) get pvc
setup:
	kubectl create namespace $(NS)
	@echo $(NS) is the namespace
	kubectl delete configmap rq-config -n $(NS) --ignore-not-found
	kubectl create configmap rq-config \
		--from-file=./templates/pod-template.json \
		-n $(NS)
	kubectl create -n $(NS) -f deploy/service_account.yaml
	kubectl create -n $(NS) -f deploy/role.yaml
	kubectl create -n $(NS) -f deploy/role_binding.yaml
	kubectl create -n $(NS) -f deploy/crds/rqcluster.example.com_rqclusters_crd.yaml
operatorimage:   
	operator-sdk build quay.io/$(IMAGEUSER)/rqlite-operator:v0.0.2
pushoperatorimage:   
	docker push quay.io/$(IMAGEUSER)/rqlite-operator:v0.0.2
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

