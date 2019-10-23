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
	export OPERATOR_NAME=rq-operator
	operator-sdk up local --namespace=rqnamespace
cleanup:   
	kubectl delete rqclusters --all -n rqnamespace
	kubectl delete deploy --all -n rqnamespace
operatorimage:   
	operator-sdk build quay.io/jemccorm/rqlite-operator:v0.0.1
