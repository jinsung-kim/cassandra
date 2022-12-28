#!/bin/bash

read -p "How many nodes per ring? " nodes
read -p "Replication factor: " rf

# for loop to create all of the instances
echo "Creating instances"

STARTING_PORT=8090
PORT_NUMBERS=()
for ((rep=0; rep < rf; rep++))
do
    for ((port=STARTING_PORT; port < STARTING_PORT + nodes; port++))
    do
        PORT_NUMBERS+=($port)
    done
    ((STARTING_PORT = STARTING_PORT + nodes))
done

BACKENDS=()
END_POINT=rf*nodes
# create corresponding backends as string
for ((i=0; i < END_POINT; i += rf))
do
    CORRESPONDING=""
    for ((j=i; j<i+rf; j++))
    do
        CORRESPONDING+=":${PORT_NUMBERS[j]},"
    done

    for ((j=i; j<i+rf; j++))
    do
        BACKENDS+=("${CORRESPONDING%,*}")
    done
done

for ((i=0; i < END_POINT; i++))
do
    go run . --listen "${PORT_NUMBERS[i]}" --backends "${BACKENDS[i]}" &
done