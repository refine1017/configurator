package setting

import (
	"github.com/go-macaron/session"
	"github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"strings"
)

var (
	// SessionConfig difines Session settings
	SessionConfig session.Options
)

func newSessionService() error {
	SessionConfig.Provider = Cfg.Section("session").Key("PROVIDER").In("memory",
		[]string{"memory", "file", "redis", "mysql", "postgres", "couchbase", "memcache", "nodb"})
	SessionConfig.ProviderConfig = strings.Trim(Cfg.Section("session").Key("PROVIDER_CONFIG").MustString(path.Join(AppDataPath, "sessions")), "\" ")
	if SessionConfig.Provider == "file" && !filepath.IsAbs(SessionConfig.ProviderConfig) {
		SessionConfig.ProviderConfig = path.Join(AppWorkPath, SessionConfig.ProviderConfig)
	}
	SessionConfig.CookieName = Cfg.Section("session").Key("COOKIE_NAME").MustString("i_like_gitea")
	SessionConfig.CookiePath = AppSubURL
	SessionConfig.Secure = Cfg.Section("session").Key("COOKIE_SECURE").MustBool(false)
	SessionConfig.Gclifetime = Cfg.Section("session").Key("GC_INTERVAL_TIME").MustInt64(86400)
	SessionConfig.Maxlifetime = Cfg.Section("session").Key("SESSION_LIFE_TIME").MustInt64(86400)

	logrus.Info("Session Service Enabled")

	return nil
}

func init() {
	registerParser(newSessionService)
}
