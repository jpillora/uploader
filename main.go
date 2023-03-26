package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jpillora/opts"
	"github.com/jpillora/requestlog"
	uploader "github.com/jpillora/uploader/lib"
)

var version = "0.0.0"

func main() {
	//cli config
	config := struct {
		Port            int  `help:"listening port"`
		NoLog           bool `help:"disable request logging"`
		uploader.Config `type:"embedded"`
	}{
		Port:   3000,
		Config: uploader.Config{},
	}

	opts.New(&config).
		Name("uploader").
		Repo("github.com/jpillora/uploader").
		Version(version).
		Parse()

	log.Printf("listening on %d...", config.Port)

	h := uploader.New(config.Config)

	if !config.NoLog {
		h = requestlog.Wrap(h)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), h)
}
