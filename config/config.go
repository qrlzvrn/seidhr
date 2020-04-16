package config

import (
	"os"
	"strconv"
)

// Используя библиотеку envconfig считываем в структуры необходимые данные из переменных окружения.
// В случае работы с докером переменные окружения прописываются в .env файле и передаются в докер.
// В случе, когда мы не хотим работать с докером, придется прописать все переменные вручную.
// Позже, может быть автоматизирую данный процесс.

// DB - конфиг для работы с базой данных
type DB struct {
	Host     string
	Port     int
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

// NewDBConf - генерирует новый конфиг для работы с базой данных
func NewDBConf() (DB, error) {
	db := DB{}
	// Приводим portк int, так как из .env считываются строки
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return db, err
	}
	db.Port = port

	db.Host = os.Getenv("DB_HOST")
	db.Username = os.Getenv("DB_USERNAME")
	db.Password = os.Getenv("DB_PASSWORD")
	db.Name = os.Getenv("DB_NAME")

	return db, nil
}

// NewTgBotConf - генерирует новый конфиг с информацией о телеграм боте
func NewTgBotConf() TgBot {
	tgBot := TgBot{}

	tgBot.APIToken = os.Getenv("TELEGRAM_APITOKEN")

	return tgBot
}

// NewSSLConf - генерирует новый конфиг с информацией о SSL
func NewSSLConf() SSL {
	ssl := SSL{}

	ssl.Fullchain = os.Getenv("FULLCHAIN")
	ssl.Privkey = os.Getenv("PRIVKEY")

	return ssl
}
