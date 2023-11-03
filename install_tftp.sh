#!/bin/bash

if ! dpkg -l | grep "tftpd-hpa" &> /dev/null; then 
    echo "Installing TFTP server (tftpd-hpa)"
    sudo apt update && sudo apt install -y tftpd-hpa
else 
    echo "TFTP server (tftpd-hpa) is already installed"
fi

if ! dpkg -l | grep "tftp-hpa" &> /dev/null; then
    echo "Installing TFTP client (tftp-hpa)"
    sudo apt update && sudo apt install -y tftp-hpa
else
    echo "TFTP client (tftp-hpa) is already installed"
fi

echo "Use: 'tftp 127.0.0.1' to connect to TFTP server running on your host machine."