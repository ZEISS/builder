package sites

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils"
)

const (
	// EmptyPrefix is the empty prefix for the key.
	EmptyPrefix = ""
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c fiber.Ctx) bool
	// Storage is the storage backend for the files.
	// This can be Azure Blob Storage, S3, or any other storage backend.
	Storage fiber.Storage
	// Prefix is a prefix for the key.
	Prefix string `json:"prefix"`
	// The value for the Cache-Control HTTP-header
	// that is set on the file response. MaxAge is defined in seconds.
	//
	// Optional. Default value 0.
	MaxAge int `json:"max_age"`
	// File to return if path is not found. Useful for SPA's.
	//
	// Optional. Default: ""
	NotFoundFile string `json:"not_found_file"`
	// The value for the Content-Type HTTP-header
	// that is set on the file response
	//
	// Optional. Default: ""
	ContentTypeCharset string `json:"content_type_charset"`
}

// ConfigDefault is the default config.
var ConfigDefault = Config{
	Next:               nil,
	Storage:            nil,
	Prefix:             "",
	MaxAge:             0,
	ContentTypeCharset: "",
}

// New creates a new middleware handler.
//
// filesystem does not handle url encoded values (for example spaces)
// on it's own. If you need that functionality, set "UnescapePath"
// in fiber.Config.
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		if cfg.NotFoundFile != "" && !strings.HasPrefix(cfg.NotFoundFile, "/") {
			cfg.NotFoundFile = "/" + cfg.NotFoundFile
		}
	}

	cacheControlStr := "public, max-age=" + strconv.Itoa(cfg.MaxAge)

	// Return new handler
	return func(c fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		method := c.Method()
		storage := cfg.Storage

		// We only serve static assets on GET or HEAD methods
		if method != fiber.MethodGet && method != fiber.MethodHead {
			return c.Next()
		}

		// Get site from domain param
		site := fiber.DomainParam(c, "site")
		path := filepath.Clean(filepath.Join(cfg.Prefix, site, c.Path()))

		if len(path) > 1 {
			path = utils.TrimRight(path, '/')
		}

		file, err := storage.GetWithContext(c.Context(), path)
		if err != nil {
			return c.Status(fiber.StatusNotFound).Next()
		}

		mimeType := http.DetectContentType(file)

		c.Status(fiber.StatusOK)
		c.Set(fiber.HeaderContentType, mimeType)
		c.Set(fiber.HeaderCacheControl, cacheControlStr)

		return c.Send(file)
	}
}
