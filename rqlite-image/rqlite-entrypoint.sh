#!/bin/bash

rqlited -http-addr 0.0.0.0:4001 -raft-addr 0.0.0.0:4002 /rqlite/file/data
#rqlited -http-addr 0.0.0.0:4001 -raft-addr 0.0.0.0:4002 -join rqcluster:4001 /rqlite/file/data
