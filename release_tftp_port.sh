#!/bin/bash

if sudo lsof -i :$1 &> /dev/null; then
	echo "Port 69 is in use. Stopping TFTP service..."
	sudo systemctl stop tftpd-hpa
	echo "TFTP service is stopped. Port 69 - released."
else
	echo "Port 69 is not in use. Free to use it."
fi
