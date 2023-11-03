package server

import (
	"bytes"
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ellofae/tftp-server/config"
	"github.com/ellofae/tftp-server/packets"
)

type Server struct {
	Retries uint8
	Timeout time.Duration
}

func getFilePayload(tftp_dir string, filename string, clientAddr string) ([]byte, error) {
	if _, err := os.Stat(tftp_dir + "/" + filename); err != nil {
		if os.IsNotExist(err) {
			log.Printf("[server] file '%s' does not exist, error: %s", filename, err.Error())
			return nil, err
		}

		log.Printf("[server] cannot get file info for %s, error: %v", filename, err)
		return nil, err
	}

	dataPayload, err := os.ReadFile(tftp_dir + "/" + filename)
	if err != nil {
		log.Printf("[server] unable to get '%s' by %s, error: %s", filename, clientAddr, err.Error())
		return nil, err
	}

	return dataPayload, nil
}

func ConfigureServer(cfg *config.Config) (*Server, error) {
	retries, err := strconv.Atoi(cfg.ServerConfiguration.Retries)
	if err != nil {
		return nil, err
	}

	if retries < 1 || retries > 255 {
		return nil, errors.New("number of retries is supposed to be a byte (1-255)")
	}

	timeout, err := strconv.Atoi(cfg.ServerConfiguration.Timeout)
	if err != nil {
		return nil, err
	}

	return &Server{
		Retries: uint8(retries),
		Timeout: time.Second * time.Duration(timeout),
	}, nil
}

func (s Server) ListenAndServe(addr string, tftp_dir string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close }()

	log.Printf("TFTP server is listening on %s ...\n", conn.LocalAddr())

	return s.Serve(conn, tftp_dir)
}

func (s *Server) Serve(conn net.PacketConn, tftp_dir string) error {
	if conn == nil {
		return errors.New("server: nil connection")
	}

	var readPacket packets.ReadRequest

	for {
		buf := make([]byte, packets.PacketSize)

		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("[server] error reading from connection %s, error %v", addr, err)
			continue
		}

		err = readPacket.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("[%s] bad read request: %v", addr, err)
			continue
		}

		// handling DATA packet sending
		go s.handleReadPacket(addr.String(), readPacket, tftp_dir)
	}
}

func (s Server) handleReadPacket(clientAddr string, readPacket packets.ReadRequest, tftp_dir string) {
	log.Printf("[%s] requested file: %s", clientAddr, readPacket.Filename)

	conn, err := net.Dial("udp", clientAddr)
	if err != nil {
		log.Printf("[server]: unable to establish dial udp connection with %s, error: %v", clientAddr, err)
		return
	}
	defer func() { _ = conn.Close() }()

	// getting data payload from TFTP server if the requested data exists
	dataPayload, err := getFilePayload(tftp_dir, readPacket.Filename, clientAddr)
	if err != nil {
		return
	}

	var (
		ackPacket  packets.AckPacket
		errPacket  packets.ErrorPacket
		dataPacket = packets.DataPacket{Payload: bytes.NewReader(dataPayload)}
		buf        = make([]byte, packets.PacketSize)
	)

NEXTPACKET:
	for n := packets.PacketSize; n == packets.PacketSize; {
		data, err := dataPacket.MarshalBinary()
		if err != nil {
			log.Printf("[server] unable to prepare data packet for %v, error: %v", clientAddr, err)
			return
		}

	RETRY:
		for i := s.Retries; i > 0; i-- {
			n, err = conn.Write(data)
			if err != nil {
				// if connection with the client cannot be established, then the DATA packet is lost
				log.Printf("[server]: unable to send packet to client on address %s, error: %v", clientAddr, err)
				return
			}

			// wait for the clien't ACK packet on recieving the data packet
			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))

			_, err = conn.Read(buf)
			if err != nil {
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					// if the ACK packet for the client hasn't been received, then send the DATA packet again
					log.Printf("[server] ACK packet from %s was not recieved because of timeout", clientAddr)
					continue RETRY
				}

				log.Printf("[server] ACK packet from %s was not received because of internal error, error: %v", clientAddr, err)
				return
			}

			switch {
			case ackPacket.UnmarshalBinary(buf) == nil:
				// ACK packet from client has been received
				if uint16(ackPacket.BlockNumber) == dataPacket.BlockNumber {
					// received ACK; send next data packet
					continue NEXTPACKET
				}
			case errPacket.UnmarshalBinary(buf) == nil:
				log.Printf("[server] an error was recevied instead of ACK packet from %s", clientAddr)
				return
			default:
				log.Printf("[server] bad DATA packet from %s", clientAddr)
			}
		}

		log.Printf("[server] exhausted retries on sending a DATA packet to %s", clientAddr)
		return
	}

	log.Printf("[server] sent %d DATA packet blocks to %s", dataPacket.BlockNumber, clientAddr)
}
