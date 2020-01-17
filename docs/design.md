# Design

## Philosophy

The rqlite-operator is meant to be lightweight in design and features but complete enough to allow it to be deployed into a production use.  The project was conceived to show how a distributed database operator could be developed using the latest frameworks such as operator-sdk as the basis of the development.

## Core Tech

The rqlite-operator is written in golang and uses the operator-sdk as its controller framework.

Both Kubernetes and OpenShift are platforms which are used for developing and testing the operator.  The operator is developed using the 1.17 version of Kubernetes but is known to work on earlier versions.

The storage (dynamic storage class) being used for development is storageos, but other storage choices are known to work as well.

## Features

The rqlite-operator includes the following features:

 * The rqlite Docker image includes rqlite v5.1.0.
 * Creates Pods that make up the rqlite nodes, one rqlite node corresponds to a single Pod, the rqlite-operator `controls` the rqlite cluster of nodes.
 * A CRD (custom resource definition) defines an `rqcluster` resource, this resource allows for the specification of an rqlite cluster of nodes.
 * A PVC (persistent volume claim) is created for each rqlite node (Pod), based on Dynamic Storage Class that is specified in a rqcluster CR (custom resource).  The PVC name is the same as the Pod name.
 * The rqlite-operator runs as 2 Deployments, one is the leader and one is the follower, this supports a more highly available controller.
 * The rqlite-operator is built using the operator-sdk.
 * The rqlite nodes (Pod) are defined as a JSON template file, allowing users to highly customize the Pods.
 * Prometheus metrics are created by the rqlite-operator by default.
 * The rqlite Pods include an anti-affinity rule which attempts to have the rqlite leader deployed onto a different Kubernetes cluster node than where the rqlite followers are deployed to.  This supports spreading the rqlite cluster nodes across different physical compute nodes.
 * The rqlite-operator is available via the Operator Lifecycle Marketplace.
 * The rqo command line client is being developed to assist the end user in creating rqlite custom resources.

## Features Being Worked On

 * The rqlite-operator deploys into and manages a single namespace currently.  In future versions, you will be able to have the rqlite-operator watch a number of different namespaces.
 * A backup and restore capability is being planned that allows you to create a CR to initiate a backup.  Today, you have to backup the rqlite manually.
