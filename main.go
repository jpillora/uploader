package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jpillora/opts"
	"github.com/jpillora/uploader/lib"
)

var VERSION = "0.0.0"

func main() {
	//cli config
	config := struct {
		Port            int `help:"listening port"`
		uploader.Config `type:"embedded"`
	}{
		Port:   3000,
		Config: uploader.Config{Dir: "."},
	}

	opts.New(&config).
		Name("uploader").
		Repo("github.com/jpillora/uploader").
		Version(VERSION).
		Parse()

	log.Printf("listening on %d...", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), uploader.New(config.Config))
}
