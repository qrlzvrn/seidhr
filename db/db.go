package db

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/qrlzvrn/seidhr/config"
)

// ConnectToDB ...
func ConnectToDB() (*sqlx.DB, error) {
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

// CreateUser ...
func CreateUser(db *sqlx.DB, tguserID int, chatID int) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tguser (id, chat_id) VALUES ($1, $2)", tguserID, chatID)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// CheckUser ...
func CheckUser(db *sqlx.DB, tguserID int) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from tguser where id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// InitMedList ...
func InitMedList(db sqlx.DB, medLines []string) error {
	tx := db.MustBegin()

	for _, med := range medLines {
		tx.MustExec("INSERT INTO medicament (title) VALUES ($1)", med)
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// CheckMed ...
func CheckMed(db sqlx.DB, medName string) (bool, error) {

	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from medicament where title=$1)", medName).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// Subscribe - создает новую подписку пользователя на лекарство
func Subscribe(db sqlx.DB, tguserID int, medicamentID int) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO subscription (tguser_id, medicament_id) VALUES ($1, $2)", tguserID, medicamentID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Unsubscribe - отменяет у пользьователя подписку на лекарство
func Unsubscribe(db *sqlx.DB, tguserID int, medicamentID int) error {
	tx := db.MustBegin()
	tx.MustExec("DELETE FROM subscription WHERE tguser_id = $1 AND medicament_id = $2", tguserID, medicamentID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// ListSubscriptions - находит все подписки пользователя и возвращает [][]string, где
// [[id title] [id title] [id title]]
func ListSubscriptions(db *sqlx.DB, tguserID int) ([][]string, error) {
	rows, err := db.Query("SELECT id, title from medicament INNER JOIN subscription on medicament.id=subscription.medicament_id WHERE subscription.tguser_id = $1", tguserID)
	if err != nil {
		return nil, err
	}

	subscriptions := [][]string{}

	for rows.Next() {
		var id int
		var title string

		rows.Scan(&id, &title)
		subscriptions = append(subscriptions, []string{strconv.Itoa(id), title})
		defer rows.Close()
	}
	return subscriptions, nil
}

// CheckSubscriptions - проверяет наличие у пользователя подписок на лекарства и,
// если у него есть хоть одна подписка, возвращает true.
func CheckSubscriptions(db *sqlx.DB, tguserID int) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM subscription WHERE tguser_id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}
