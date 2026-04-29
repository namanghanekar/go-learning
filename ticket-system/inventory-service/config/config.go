package config

type Config struct {
	DB_DSN string
}

func Load() *Config {
	return &Config{
		DB_DSN: "root:naman@tcp(127.0.0.1:3306)/ticketdb?charset=utf8mb4&parseTime=True&loc=Local",
	}
}
