package env

type Config struct {
	SQLiteDBPath string `env:"SQLITE_DB_PATH"`
}
