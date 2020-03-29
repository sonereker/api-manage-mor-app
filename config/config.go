package config

type Config struct {
	DB  *DBConfig
	App *AppConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Dialect  string
	Username string
	Password string
	Name     string
	SSLMode  string
}

type AppConfig struct {
	Host string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Host:     "localhost",
			Port:     "5432",
			Dialect:  "postgres",
			Username: "postgres",
			Password: "",
			Name:     "kule_app_dev",
			SSLMode:  "disable",
		},
		App: &AppConfig{Host: "localhost:4000"},
	}
}
