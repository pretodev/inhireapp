package env

type Config struct {
	SQLiteDBPath            string `env:"SQLITE_DB_PATH, required"`
	ChromiumHeadlessEnabled string `env:"CHROMIUM_HEADLESS_ENABLED, default=YES"`
}

func (cfg *Config) IsChromiumHeadlessEnabled() bool {
	return cfg.ChromiumHeadlessEnabled == "YES"
}
