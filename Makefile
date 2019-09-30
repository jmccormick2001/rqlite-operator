rqlite-image:   
	sudo --preserve-env buildah bud -f ./docker/Dockerfile -t jmccormick2001/rqlite:v0.0.1 ./docker
	sudo --preserve-env buildah push jmccormick2001/rqlite:v0.0.1 docker-daemon:jmccormick2001/rqlite:v0.0.1
	docker tag docker.io/jmccormick2001/rqlite:v0.0.1  jmccormick2001/rqlite:v0.0.1
