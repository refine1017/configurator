package setting

import (
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Cache represents cache settings
type Cache struct {
	Adapter  string
	Interval int
	Conn     string
	TTL      time.Duration
}

var (
	// CacheService the global cache
	CacheService *Cache
)

func newCacheService() error {
	sec := Cfg.Section("cache")
	CacheService = &Cache{
		Adapter: sec.Key("ADAPTER").In("memory", []string{"memory", "redis", "memcache"}),
	}
	switch CacheService.Adapter {
	case "memory":
		CacheService.Interval = sec.Key("INTERVAL").MustInt(60)
	case "redis", "memcache":
		CacheService.Conn = strings.Trim(sec.Key("HOST").String(), "\" ")
	default:
		logrus.Fatalf("Unknown cache adapter: %s", CacheService.Adapter)
	}
	CacheService.TTL = sec.Key("ITEM_TTL").MustDuration(16 * time.Hour)

	logrus.Info("Cache Service Enabled")

	return nil
}

func init() {
	registerParser(newCacheService)
}
