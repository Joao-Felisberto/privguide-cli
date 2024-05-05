package util

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
)

const (
	RESET = "\033[0m"

	BLACK        = 30
	RED          = 31
	GREEN        = 32
	YELLOW       = 33
	BLUE         = 34
	MAGENTA      = 35
	CYAN         = 36
	LIGHTGRAY    = 37
	DARKGRAY     = 90
	LIGHTRED     = 91
	LIGHTGREEN   = 92
	LIGHTYELLOW  = 93
	LIGHTBLUE    = 94
	LIGHTMAGENTA = 95
	LIGHTCYAN    = 96
	WHITE        = 97
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, RESET)
}

type HumanFriendlyHandler struct {
	handler    slog.Handler
	byteBuffer *bytes.Buffer
	m          *sync.Mutex
}

func suppressLogDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func NewHumanFriendlyHandler(opts *slog.HandlerOptions) slog.Handler {
	/*
		if opts == nil {
			opts = &slog.HandlerOptions{}
		}
		return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressLogDefaults(opts.ReplaceAttr),
		})
	*/

	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &HumanFriendlyHandler{
		byteBuffer: b,
		handler: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressLogDefaults(opts.ReplaceAttr),
		}),
		m: &sync.Mutex{},
	}
}

func (h *HumanFriendlyHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.handler.Enabled(ctx, lvl)
}

func (h *HumanFriendlyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *HumanFriendlyHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

func (h *HumanFriendlyHandler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String() + ": "

	switch record.Level {
	case slog.LevelDebug:
		level = colorize(DARKGRAY, level)
	case slog.LevelInfo:
		level = colorize(CYAN, level)
	case slog.LevelWarn:
		level = colorize(LIGHTYELLOW, level)
	case slog.LevelError:
		level = colorize(LIGHTRED, level)
	}

	fmt.Print(
		level,
		colorize(WHITE, record.Message),
	)
	record.Attrs(func(a slog.Attr) bool {
		switch t := a.Value.Any().(type) {
		case []uint8:
			fmt.Printf(" %s: '%+v'", colorize(LIGHTGRAY, a.Key), colorize(CYAN, string(t)))
		default:
			fmt.Printf(" %s: '%s'", colorize(LIGHTGRAY, a.Key), colorize(CYAN, fmt.Sprintf("%+v", t)))
		}
		return true
	})
	fmt.Print("\n")

	return nil
}
