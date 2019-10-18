a Kube operator for the rqlite database.

### How to Run It

    mkdir myproject; cd myproject
    git clone https://github.com/jmccormick2001/rqlite-operator
    cd rqlite-operator

As a cluster-admin run the following:

    kubectl create namespace rqnamespace
    make setup-as-cluster-admin

As a normal user (or cluster-admin) run the following:

    make setup

Then to start up the rqlite-operator and a test rqlite cluster:

    make testit

### How to Test It

Verify that the rqlite pods are running

    kubectl -n rqnamespace get pods

Run some sample SQL commands inside one of the pods (e.g. example-rqcluster-khig):

    kubectl -n rqnamespace exec -it example-rqcluster-khig bash
    rqlite -H example-rqcluster-khig
    create table mytesttable (id int);
    insert into mytesttable values (2);
    select * from mytesttable;
    exit
    exit

Then to verify that the rqlite cluster is replicating state, exec into a differenent rqlite cluster pod (e.g. example-rqcluster-lxvy):

    kubectl -n rqnamespace exec -it example-rqcluster-lxvy bash
    rqlite -H example-rqcluster-lxvy
    select * from mytesttable;
   
If all is working, you should see the same inserted data in the
other rqlite Pods. 
