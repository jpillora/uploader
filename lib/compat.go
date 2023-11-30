package uploader

import (
	"net/http"

	"github.com/jpillora/uploader/internal/handler"
)

type Config = handler.Config

func New(c Config) http.Handler {
	h, err := handler.New(c)
	if err != nil {
		panic(err)
	}
	return h
}
