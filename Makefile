.PHONY: release-tftp-port

release-tftp-port:
	bash release_tftp_port.sh 69

.PHONY: run-server

run-server: release-tftp-port
	GOOS=linux GOARCH=amd64 go build -o main main.go
	sudo ./main

.PHONY: install-tftp

install-tftp:
	bash install_tftp.sh