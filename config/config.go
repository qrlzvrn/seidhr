package config

import (
	"github.com/BurntSushi/toml"
)

// Config - основной конфиг, содержащий в себе все остальные
type Config struct {
	DB    DB
	TgBot TgBot
	SSL   SSL
}

// DB - конфиг для работы с базой данных
type DB struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

// TgBot - конфиг для работы с телеграм ботом, пока что нужен только для хранения токена, в будущем может быть расширен
type TgBot struct {
	APIToken string
}

// SSL - конфиг хранящий пути к сертификатам
type SSL struct {
	Fullchain string
	Privkey   string
}

// InitConf - читает конфигурационный файл и заполняет структуру Config, после чего возвращает ее
func InitConf() (Config, error) {
	var conf Config

	if _, err := toml.DecodeFile("/config.toml", &conf); err != nil {
		return conf, err
	}
	return conf, nil
}
