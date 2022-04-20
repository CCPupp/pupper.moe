#!/bin/bash
until ./server; do
	echo "Server Crashed! Restarting..."
  	sleep 1
done
