#!/bin/bash

#FILE="rel-simple-beam-form.go"
FILE="rel-prob1to5"

while true; do
    inotifywait -q -e modify $FILE
    go run $FILE
done