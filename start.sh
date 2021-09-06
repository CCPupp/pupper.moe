#!/bin/bash
while ! ./server
do
  sleep 1
  echo "Restarting program..."
done
