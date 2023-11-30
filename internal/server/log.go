package server

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type shandler struct {
	group   string
	verbose bool
	attrs   []slog.Attr
}

func (h *shandler) Enabled(_ context.Context, l slog.Level) bool {
	if h.verbose {
		return true
	}
	return l >= slog.LevelInfo
}

func (h *shandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := *h
	h2.attrs = append(h.attrs, attrs...)
	return &h2
}

func (h *shandler) WithGroup(name string) slog.Handler {
	h2 := *h
	if h2.group == "" {
		h2.group = name
	} else {
		h2.group += "." + name
	}
	return &h2
}

func (h *shandler) Handle(ctx context.Context, r slog.Record) error {
	sb := strings.Builder{}
	sb.WriteString(r.Time.Format(time.RFC3339))

	sb.WriteRune(' ')
	sb.WriteString(r.Level.String())

	if h.group != "" {
		sb.WriteRune(' ')
		sb.WriteString(h.group)
	}

	sb.WriteRune(' ')
	sb.WriteString(r.Message)

	add := func(attr slog.Attr) bool {
		sb.WriteRune(' ')
		sb.WriteString(attr.Key)
		sb.WriteRune('=')
		sb.WriteString(attr.Value.String())
		return true
	}
	r.Attrs(add)
	for _, attr := range h.attrs {
		add(attr)
	}
	fmt.Println(sb.String())
	return nil
}
