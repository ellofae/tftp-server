package main

import (
	"log"

	"github.com/ellofae/tftp-server/config"
	"github.com/ellofae/tftp-server/server"
)

func main() {
	cfg := config.ParseConfig(config.ConfigureViper())

	srv, err := server.ConfigureServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tftpServerAddress := cfg.ServerConfiguration.Address
	tftpServerDirectory := cfg.ServerConfiguration.TFTP_directory

	log.Fatal(srv.ListenAndServe(tftpServerAddress, tftpServerDirectory))
}
