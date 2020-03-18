package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/qrlzvrn/seidhr/config"
)

func ConnectToBD() (*sqlx.DB, error) {
	dbConf, err := config.NewDBConf()
	if err != nil {
		return nil, err
	}
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.Name)

	if db, err := sqlx.Connect("postgres", dbInfo); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
