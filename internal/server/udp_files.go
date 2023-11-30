package server

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jpillora/sizestr"
	"github.com/jpillora/uploader/internal/x"
)

type udpFiles struct {
	log      *slog.Logger
	config   Config
	mut      sync.Mutex
	files    map[string]*udpFile
	sweeping bool
}

func (cs *udpFiles) upsert(id string) (*udpFile, error) {
	cs.mut.Lock()
	defer cs.mut.Unlock()
	if cs.files == nil {
		cs.files = map[string]*udpFile{}
	}
	f, ok := cs.files[id]
	if !ok {
		// open file
		path := x.Join(cs.config.Overwrite, cs.config.Dir, id+".bin")
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
		// prepare new file
		f = &udpFile{
			log:  cs.log.WithGroup(id),
			id:   id,
			file: file,
			path: path,
		}
		cs.log.Info("receiving", "path", path)
		cs.files[id] = f
	}
	return f, nil
}

func (cs *udpFiles) enqueueSweep() {
	cs.mut.Lock()
	defer cs.mut.Unlock()
	if cs.sweeping {
		return
	}
	cs.sweeping = true
	time.AfterFunc(cs.config.UDPClose, cs.sweep)
	cs.log.Debug("enque-sweep")
}

func (cs *udpFiles) sweep() {
	cs.mut.Lock()
	defer cs.mut.Unlock()
	closed := 0
	for id, c := range cs.files {
		if c.timeout() {
			c.close()
			delete(cs.files, id)
			closed++
		}
	}
	cs.log.Debug("swept", "closed", closed, "remaining", len(cs.files))
	cs.sweeping = false
}

func (cs *udpFiles) del(id string) {
	cs.mut.Lock()
	defer cs.mut.Unlock()
	if cs.files == nil {
		return
	}
	c, ok := cs.files[id]
	if ok {
		c.del()
	}
	delete(cs.files, id)
	cs.log.Info("delete", "id", id)
}

type udpFile struct {
	log   *slog.Logger
	id    string
	path  string
	size  int64
	mut   sync.Mutex
	file  *os.File
	mtime time.Time
}

func (f *udpFile) write(p []byte) (n int, err error) {
	f.mut.Lock()
	defer f.mut.Unlock()
	n, err = f.file.Write(p)
	f.size += int64(n)
	return n, err
}

func (f *udpFile) timeout() bool {
	f.mut.Lock()
	defer f.mut.Unlock()
	return time.Since(f.mtime) > 5*time.Second
}

func (f *udpFile) close() {
	f.mut.Lock()
	defer f.mut.Unlock()
	if f.file != nil {
		f.log.Info("received", "size", sizestr.ToString(f.size))
		f.file.Close()
		f.file = nil
	}
}

func (f *udpFile) del() {
	f.close()
	os.Remove(f.path)
	f.log.Info("deleted")
}
