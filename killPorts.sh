#!/bin/bash

read -p "How many nodes per ring? " nodes
read -p "Replication factor: " rf

echo "Ending all ports"

STARTING_PORT=8090
PORT_NUMBERS=()
for ((rep=0; rep < rf; rep++))
do
    for ((port=STARTING_PORT; port < STARTING_PORT + nodes; port++))
    do
        kill $( lsof -i:${port} -t ) &
    done
    ((STARTING_PORT = STARTING_PORT + nodes))
done