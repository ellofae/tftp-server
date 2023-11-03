#!/bin/bash

if sudo lsof -i :$1 &> /dev/null; then
	echo "Port $1 is in use. Stopping TFTP service..."
	sudo service tftpd-hpa stop
	echo "TFTP service is stopped. Port $1 - released."
else
	echo "Port $1 is not in use. Free to use it."
fi
