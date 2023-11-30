package server

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jpillora/ipfilter"
	"github.com/jpillora/requestlog"
	"github.com/jpillora/uploader/internal/handler"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Dir       string        `help:"output directory (defaults to tmp)"`
	Overwrite bool          `help:"duplicates are overwritten (auto-renames files by default)"`
	Auth      string        `help:"require basic auth 'username:password' on http connections"`
	Port      int           `help:"tcp listening port"`
	UDPPort   int           `help:"udp listening port (default disabled)"`
	UDPClose  time.Duration `help:"close udp file after timeout"`
	NoLog     bool          `help:"disable http request logging"`
	AllowedIP []string      `opts:"short=i, help=allowed ip range"`
	Verbose   bool          `help:"enable verbose logging"`
}

type Server struct {
	Config
	filter *ipfilter.IPFilter
	log    *slog.Logger
}

func New(c Config) *Server {
	if c.Dir == "" {
		c.Dir = os.TempDir()
	}

	log := slog.New(&shandler{
		verbose: c.Verbose,
	}).WithGroup("uploader")

	return &Server{
		Config: c,
		filter: nil,
		log:    log,
	}
}

func (s *Server) Start() error {
	if len(s.AllowedIP) > 0 {
		s.log.Info("enable filter", "allowed-ips", s.AllowedIP)
		s.filter = ipfilter.New(ipfilter.Options{
			AllowedIPs:     s.AllowedIP,
			BlockByDefault: true,
		})
	}

	if info, err := os.Stat(s.Dir); err != nil || !info.IsDir() {
		return fmt.Errorf("invalid directory: %s", s.Dir)
	}
	s.log.Info("output directory", "path", s.Dir)

	eg := errgroup.Group{}
	if s.UDPPort != 0 {
		eg.Go(s.startUDP)
	}
	eg.Go(s.startHTTP)
	return eg.Wait()
}

func (s *Server) startHTTP() error {
	log := s.log.WithGroup("http")
	// for compatibility with the old version, we keep the handler config separate
	h, err := handler.New(handler.Config{
		Dir:       s.Dir,
		Overwrite: s.Overwrite,
		Auth:      s.Auth,
		Logger:    log,
	})
	if err != nil {
		return err
	}
	if s.filter != nil {
		h = s.filter.Wrap(h)
	}
	if !s.NoLog {
		// TODO: add slog support to requestlog
		h = requestlog.Wrap(h)
	}
	addr := fmt.Sprintf("0.0.0.0:%d", s.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Info("listening on tcp", "addr", addr)
	return http.Serve(l, h)
}

func (s *Server) startUDP() error {
	log := s.log.WithGroup("udp")
	addr := fmt.Sprintf("0.0.0.0:%d", s.UDPPort)
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", uaddr)
	if err != nil {
		return fmt.Errorf("udp: %w", err)
	}
	log.Info("listening on udp", "addr", addr)
	// track udp files
	files := &udpFiles{
		log:    log,
		config: s.Config,
	}
	// read all udp packets
	// TODO: writer-pool for each file to allow concurrent writes
	buff := make([]byte, 32*1024)
	for {
		n, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			log.Warn("read", "err", err)
			continue
		}
		// ip check
		if s.filter != nil && !s.filter.NetAllowed(addr.IP) {
			log.Debug("blocked", "addr", addr)
			continue
		}
		// prepare id
		hash := md5.Sum([]byte(addr.String()))
		id := hex.EncodeToString(hash[:])[0:8]
		// prepare file
		f, err := files.upsert(id)
		if err != nil {
			log.Warn("get-err", "id", id, "err", err)
			continue
		}
		// write to file
		b := buff[:n]
		if _, err := f.write(b); err != nil {
			log.Warn("write-err", "id", id, "err", err)
			files.del(id)
		}
		// sweep all files soon
		files.enqueueSweep()
	}
}
