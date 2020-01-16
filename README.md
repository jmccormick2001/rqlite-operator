---
description: A Kubernetes operator for the rqlite database.
---

# rqlite-operator

## Getting Started

You can install the rqlite-operator as follows:

```
$ make setup
$ make install
$ make test
```

{% hint style="info" %}
 This example above assumes you have a working k8s cluster available and you are running as a cluster-admin user.
{% endhint %}

Verify the rqlite-operator is running and a sample rqlite deployment of 3 nodes is running:

```text
$ kubectl -n rq get deployment
$ kubectl -n rq get pod
```



