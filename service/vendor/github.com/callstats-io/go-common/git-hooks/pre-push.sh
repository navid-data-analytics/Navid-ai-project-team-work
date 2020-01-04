#!/bin/bash

echo 'Running tests...'
docker-compose run test_go_common
status=$?
if [ $status != 0 ]; then exit $status; fi
