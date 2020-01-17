---
description: Instructions if you want to build the rqlite-operator from source.
---

# Building from Source

## Prerequisites

For this build, I am assuming you are working on a Linux
system that has the following packages installed:
 
 * buildah
 * docker
 * golang 1.13+

Modify the Makefile environment variables to your local environment
by changing the following in the Makefile:

 * NS - this is the namespace into which the rqlite-operator will be deployed,
`rq` is the default value 
 * IMAGEUSER - this is the user name for the rqlite-operator and rqlite images which are part of the image path, `someuser` is the default value
 
## Build Steps

Building the rqlite-operator from source code includes
the following:

```bash
mkdir rqlite-operator
git clone https://github.com/jmccormick2001/rqlite-operator.git
cd rqlite-operator
make rqliteimage
make operatorimage
```
