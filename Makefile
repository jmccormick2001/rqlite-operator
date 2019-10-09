rqliteimage:   
	sudo --preserve-env buildah bud -f ./rqlite-image/Dockerfile -t jmccormick2001/rqlite:v0.0.1 ./rqlite-image
	sudo --preserve-env buildah push jmccormick2001/rqlite:v0.0.1 docker-daemon:jmccormick2001/rqlite:v0.0.1
	docker tag docker.io/jmccormick2001/rqlite:v0.0.1  jmccormick2001/rqlite:v0.0.1
configmap:   
	kubectl delete configmap rq-config
	kubectl create configmap rq-config --from-file=./templates
testit:   
	kubectl create -f deploy/operator.yaml
	sleep 10
	kubectl create -f deploy/crds/rqcluster.example.com_v1alpha1_rqcluster_cr.yaml
cleanup:   
	kubectl delete rqclusters --all
	kubectl delete deploy --all
