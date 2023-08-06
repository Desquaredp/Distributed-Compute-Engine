#!/bin/bash

for i in {01..12}; do
  ssh orion$i 'cd P1-pickle-rick/src; go run -race spawn/spawn.go' &
done

echo "Spawned processes in the background."