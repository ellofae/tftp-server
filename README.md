# tftp-server

Trivial File Transfer Protocol (TFTP) is a simple lockstep File Transfer Protocol which allows a client to get a file from or put a file onto a remote host. One of its primary uses is in the early stages of nodes booting from a local area network. TFTP has been used for this application because it is very simple to implement.

Due to its simple design, TFTP can be easily implemented by code with a small memory footprint. It is therefore the protocol of choice for the initial stages of any network booting strategy like BOOTP, PXE, BSDP, etc.

Due to TFTP DATA and ACK packets 2-bytes Block Counter, the maximum file size that the server is able to send to the client is `~= 33.8 MB` (516 bytes * 65535)

## TFTP protocol packets

TFTP operates using various packet types that serve different functions in the file transfer process.

### READ packet (Read Request):

The primary purpose of the RRQ packet is to request a specific file from the TFTP server. The client sends this packet to the server to initiate the process of reading a file. The client specifies the filename it wants to retrieve and the mode of transfer (octet, netascii).

    RRQ packet is composed of Operation Code(2-bytes), Filename(n-bytes), 1-byte delimiter, Mode(n-bytes) and 1-byte delimiter.

### DATA packet:

TFTP divides the file into smaller blocks or packets for efficient transfer. The DATA packet is responsible for carrying a portion of the file's data.

The size of each DATA packet is relatively small (typically 516 bytes): 512 bytes for file payload and 4 bytes for packet header.

    DATA packet is composed of Operation Code(2-bytes), Package Number(2-bytes) and Payload(up to 512 bytes).

### ACK packet (Acknowledgment):

After the client receives a data packet, it sends an ACK packet to the server. The ACK packet essentially confirms that the data packet was successfully received without errors. This acknowledgment is crucial because it informs the server that it can send the next data packet.

    ACK packet is composed of Operation Code(2-bytes) and Recieved Packet Number(2-bytes).

### ERROR packet:

The ERR packet is used to notify the sender or receiver of a TFTP operation about an error that has occurred during the transfer.

    RFC 1350 Error codes:
    
    0: Not defined, see error message.
    1: File not found.
    2: Access violation.
    3: Disk full or allocation exceeded.
    4: Illegal TFTP operation.
    5: Unknown transfer ID.
    6: File already exists.
    7: No such user.
    8: Requested file is not valid for this server.

## Data - Acknowledgment packets workflow

The server requires an acknowledgment from the client after each data
packet. If the server does not receive a timely acknowledgment or an error
from the client, the server will retry the transmission until it receives a reply
or exhausts its number of retries.

![main](https://i.imgur.com/qiFnofN.png)

## Ways to install and launch TFTP server

Before you proceed with running the TFTP server on Linux, make sure you have installed the TFTP server and client on your host machine.

### Installation
In order to install it, you need to run `make install-tftp port` or run the `./install_tftp.sh port` bash script. TFTP server and client will be installed.

By default, TFTP server is running on port **:69**, which menas, if you intend to run the tftp-server on this port, you need to release it.


### Port release and client-server connection
To release the used **:69** port, run `make release-tftp-port` or run the `./release_tftp_port.sh` bash script.

Then connect the with TFTP client to the TFTP server with the help of command: `tftp 127.0.0.1`, which will connect the TFTP client with the TFTP server running on your host machine.

If the TFTP server is run not on the local machine, then use different address for connection.


### Run the server
To run the TFTP server, you can run `make run-server` or do it step-by-step:

1) Build binary file: GOOS=linux GOARCH=amd64 go build -o main main.go

2) Run the binary file with root privileges if you intend to use server privileged ports below 1024


## Configuration

The TFTP server uses .yaml file to configure. Here is the example of a config.yaml file used for configuration:

    ServerConfiguration:
    retries: '6'
    timeout: '10'
    address: '127.0.0.1:69'
    tftp_directory: '/srv/tftp'

## Usage example

TFTP server is running on port **:69** and TFTP client is connected to the server running on localhost.

![main](https://i.imgur.com/QoDenCW.png)
