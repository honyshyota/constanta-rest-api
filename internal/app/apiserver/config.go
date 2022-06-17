package apiserver

type Config struct {
	BindAddr         string `toml:"bind_addr"`
	LogLevel         string `toml:"log_level"`
	DatabaseURL      string `toml:"database_url"`
	SessionKey       string `toml:"session_key"`
	PgMigrationsPath string `toml:"pg_migrations_path"`
	PgURL            string `toml:"pg_url"`
	PgTest           string `toml:"pg_test"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		LogLevel: "debug",
	}
}
