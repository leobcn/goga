#!/bin/bash

FILES="evolver.go island.go operators.go ops_floats.go t_floats_test.go"
TEST="flt04"

refresh(){
    go test -test.run="$TEST"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
