package uploader

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/jpillora/sizestr"
)

//go:embed static/*
var content embed.FS

type Config struct {
	Dir       string `help:"output directory (defaults to tmp)"`
	Overwrite bool   `help:"duplicates are overwritten (auto-renames files by default)"`
	Auth      string `help:"require basic auth 'username:password'"`
}

func New(config Config) http.Handler {

	if config.Dir == "" {
		config.Dir = os.TempDir()
	}
	if info, err := os.Stat(config.Dir); err != nil || !info.IsDir() {
		log.Fatalf("Invalid directory: %s", config.Dir)
	}
	log.Printf("saving files to: %s", config.Dir)

	var contentHandler = gziphandler.GzipHandler(
		http.FileServer(http.FS(content)),
	)

	uploadID := 1
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

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

		path := filepath.Join(config.Dir, filename)
		if !config.Overwrite {
			count := 1
			for {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					break
				}
				ext := filepath.Ext(path)
				prefix := strings.TrimSuffix(path, ext)
				if count > 1 {
					prefix = strings.TrimSuffix(prefix, fmt.Sprintf("-%d", count))
				}
				count++
				path = fmt.Sprintf("%s-%d%s", prefix, count, ext)
			}
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Printf("file error: %s", err)
			w.WriteHeader(500)
			return
		}
		defer file.Close()

		log.Printf("#%04d receiving %s", uploadID, path)
		if n, err := io.Copy(file, part); err != nil {
			log.Printf("#%04d receive error: %s", uploadID, err)
		} else {
			log.Printf("#%04d received %s", uploadID, sizestr.ToString(n))
		}
		uploadID++
	})
}
