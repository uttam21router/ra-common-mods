package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/routerarchitects/ra-common-mods/logging"

	"github.com/routerarchitects/ra-common-mods/ctxmeta"
)

type RequestLoggerOptions struct {
	// Logger to use. If nil, middleware uses logging.Subsystem("http").
	Logger logging.Logger

	// Header to read/write request id.
	// Default: "X-Request-Id"
	RequestIDHeader string

	// If true, sets the request id header in the response.
	// Default: true
	EchoRequestID bool

	// Generator used when request id is missing.
	// Default: crypto-rand 16 bytes hex.
	GenerateRequestID func(context.Context) string

	// If true, logs both "request started" and "request completed".
	// Default: true
	LogStart bool
	LogEnd   bool
}

func RequestLogger(opts RequestLoggerOptions) fiber.Handler {
	log := opts.Logger
	if log == nil {
		log = logging.Subsystem("http")
	}

	hdr := opts.RequestIDHeader
	if hdr == "" {
		hdr = "X-Request-Id"
	}

	gen := opts.GenerateRequestID
	if gen == nil {
		gen = func(context.Context) string { return randomHex(16) }
	}

	// defaults
	echo := opts.EchoRequestID
	if !opts.EchoRequestID && opts.RequestIDHeader == "" {
		// if user didn't set options at all, echo should be true by default
		echo = true
	}
	logStart := opts.LogStart
	logEnd := opts.LogEnd
	if !opts.LogStart && !opts.LogEnd {
		// if user didn't set either explicitly, enable both
		logStart, logEnd = true, true
	}

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Fiber carries a context.Context via UserContext().
		uctx := c.UserContext()
		if uctx == nil {
			uctx = context.Background()
		}

		// 1) Determine request_id: prefer ctxmeta already in context, else header, else generate.
		rid, _ := ctxmeta.RequestID(uctx)
		if rid == "" {
			rid = strings.TrimSpace(c.Get(hdr))
		}
		if rid == "" {
			rid = gen(uctx)
		}

		// 2) Put it into user context and propagate into request headers for downstream.
		uctx = ctxmeta.WithRequestID(uctx, rid)
		c.SetUserContext(uctx)

		// Ensure header exists on request (useful for internal client calls that reuse incoming headers).
		c.Request().Header.Set(hdr, rid)

		// 3) Echo to response header (optional).
		if echo {
			c.Set(hdr, rid)
		}

		// 4) Log start.
		if logStart {
			log.Info(uctx, "incoming request",
				logging.Fields{
					"method":     c.Method(),
					"path":       c.Path(),
					"host":       c.Hostname(),
					"remote_ip":  c.IP(),
					"user_agent": c.Get("User-Agent")},
			)
		}

		// 5) Call next.
		err := c.Next()

		// 6) Log completion.
		if logEnd {
			status := c.Response().StatusCode()

			// Best-effort bytes:
			// - ContentLength may be -1 if unknown.
			// - Body length is available for many responses but may be 0 for streamed.
			bytes := c.Response().Header.ContentLength()
			if bytes < 0 {
				bytes = len(c.Response().Body())
			}

			if err != nil {
				// If you have Errorw(...) use that; otherwise log as info with error field.
				log.Info(uctx, "request completed with error",
					logging.Fields{
						"method":      c.Method(),
						"path":        c.Path(),
						"status":      status,
						"bytes":       bytes,
						"duration_ms": time.Since(start).Milliseconds(),
						"error":       err.Error()},
				)
			} else {
				log.Info(uctx, "request completed",
					logging.Fields{
						"method":      c.Method(),
						"path":        c.Path(),
						"status":      status,
						"bytes":       bytes,
						"duration_ms": time.Since(start).Milliseconds(),
					},
				)
			}
		}

		return err
	}
}

func randomHex(nbytes int) string {
	b := make([]byte, nbytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
