rqliteimage:   
	sudo --preserve-env buildah bud -f ./rqlite-image/Dockerfile -t jmccormick2001/rqlite:v0.0.1 ./rqlite-image
	sudo --preserve-env buildah push jmccormick2001/rqlite:v0.0.1 docker-daemon:jmccormick2001/rqlite:v0.0.1
	docker tag docker.io/jmccormick2001/rqlite:v0.0.1  jmccormick2001/rqlite:v0.0.1
configmap:   
	kubectl delete configmap rq-config -n rqnamespace --ignore-not-found
	kubectl create configmap rq-config --from-file=./templates -n rqnamespace
testit:   
	kubectl create -f deploy/service_account.yaml -n rqnamespace
	kubectl create -f deploy/role.yaml -n rqnamespace
	kubectl create -f deploy/role_binding.yaml -n rqnamespace
	kubectl create -f deploy/operator.yaml -n rqnamespace
	sleep 10
	kubectl create -f deploy/crds/rqcluster.example.com_v1alpha1_rqcluster_cr.yaml -n rqnamespace
cleanup:   
	kubectl delete rqclusters --all -n rqnamespace
	kubectl delete deploy --all -n rqnamespace
operatorimage:   
	operator-sdk build quay.io/jemccorm/rq-operator:v0.0.1
