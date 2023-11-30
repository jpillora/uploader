package main

import (
	"log"
	"time"

	"github.com/jpillora/opts"
	"github.com/jpillora/uploader/internal/server"
)

var version = "0.0.0"

func main() {
	config := server.Config{
		Port:     3000,
		UDPClose: 2 * time.Second,
	}
	opts.New(&config).
		Name("uploader").
		Summary("note: udp creates files using a stream of packets. udp packets are not authenticated,\n" +
			"so it's highly recommended that you set an allowed-ip range to prevent misuse.\n" +
			"udp packets are all appended to a file called 'md5(<src-ip>:<src-port>).bin'.\n" +
			"udp streams are considered closed after --udp-close and the file will be closed.").
		Repo("github.com/jpillora/uploader").
		Version(version).
		Parse()
	h := server.New(config)
	if err := h.Start(); err != nil {
		log.Fatal(err)
	}
}
