#!/bin/bash


function trap_sigterm() {
killall rqlited
}

trap 'trap_sigterm' SIGINT SIGTERM

echo $JOIN_ADDRESS is join address

rqlited -http-addr $HOSTNAME:4001 -raft-addr $HOSTNAME:4002 $JOIN_ADDRESS /rqlite/file/data
