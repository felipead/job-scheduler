#!/bin/bash

FILE=bin/job-scheduler

if [[ -x "$FILE" ]]
then
    bin/job-scheduler $@
else
    echo "Could not find or execute '$FILE'; please run \`make build\`"
fi
