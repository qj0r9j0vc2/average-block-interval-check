i#!/bin/bash

height=$1
iterations=$2
total_diff=0

for ((i=0; i<iterations; i++)); do
    block=$(curl -s "http://localhost:26657/block?height=$height" | jq -r '.result.block.header.time')
    prev_block=$(curl -s "http://localhost:26657/block?height=$((height - 1))" | jq -r '.result.block.header.time')
    diff=$(echo "scale=4; ($(date -d"$block" +"%s") - $(date -d"$prev_block" +"%s")) / 1" | bc)
    total_diff=$(echo "scale=4; $total_diff + $diff" | bc)
    height=$((height - 1))
done

average=$(echo "scale=4; $total_diff / $iterations" | bc)

echo "Average time difference: $average seconds"
