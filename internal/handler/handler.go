package handler

import (
	"embed"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/jpillora/sizestr"
	"github.com/jpillora/uploader/internal/x"
)

//go:embed static/*
var content embed.FS

type Config struct {
	Dir       string
	Overwrite bool
	Auth      string
	Logger    *slog.Logger
}

func New(config Config) (http.Handler, error) {

	log := config.Logger
	if log == nil {
		log = slog.Default()
	}
	if config.Dir == "" {
		config.Dir = os.TempDir()
	}
	if info, err := os.Stat(config.Dir); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("invalid directory: %s", config.Dir)
	}

	var contentHandler = gziphandler.GzipHandler(
		http.FileServer(http.FS(content)),
	)

	uploadID := 1
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		l := log.WithGroup(fmt.Sprintf("f%d", uploadID))
		if config.Auth != "" {
			u, p, _ := r.BasicAuth()
			if config.Auth != u+":"+p {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Access Denied"))
				return
			}
		}

		if r.Method == "GET" {
			r.URL.Path = "/static" + r.URL.Path
			contentHandler.ServeHTTP(w, r)
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(400)
			fmt.Fprint(w, "Expecting POST")
			return
		}

		multi, err := r.MultipartReader()
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprint(w, "Expecting multipart form")
			return
		}

		part, err := multi.NextPart()
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Expecting multipart form (%s)", err)
			return
		}
		defer part.Close()

		if part.FormName() != "file" {
			w.WriteHeader(400)
			fmt.Fprint(w, "Expecting multipart entry 'file'")
			return
		}

		filename := part.FileName()
		if filename == "" {
			filename = "file"
		}

		path := x.Join(config.Overwrite, config.Dir, filename)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			l.Error("open-file", "err", err)
			w.WriteHeader(500)
			return
		}
		defer file.Close()

		l.Info("receiving", "path", path)
		if n, err := io.Copy(file, part); err != nil {
			l.Error("receive-copy", "err", err)
		} else {
			l.Info("received", "size", sizestr.ToString(n))
		}
		uploadID++
	}), nil
}
