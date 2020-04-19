package med

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/qrlzvrn/seidhr/db"
	"github.com/qrlzvrn/seidhr/keyboards"
)

// Apothecary ...
type Apothecary struct {
	Name         string `json:"name"`
	Addr         string `json:"address"`
	FedExemption string `json:"ost1"`
	RegExemption string `json:"ost2"`
	PsyExemption string `json:"ost3"`
	VZN          string `json:"ost4"`
	Date         string `json:"date"`
}

//District ...
type District struct {
	Name         string       `json:"name"`
	ID           string       `json:"id"`
	Apothecaries []Apothecary `json:"apothecaries"`
}

// Result ...
type Result struct {
	Name      string     `json:"name"`
	Districts []District `json:"districts"`
}

// Model ...
type Model struct {
	Result []Result `json:"result"`
}

// Jsn - основная структура ответа
type Jsn struct {
	Status string `json:"status"`
	Model  Model  `json:"model,omitempty"`
	Errors string `json:"errors,omitempty"`
}

// ReqMedInfo - опрашивает сторонний сервис о наличии лекарства,
// анмаршалит полученный json в структуру Jsn и возвращает ее
func ReqMedInfo(medTitle string) (Jsn, error) {
	hh := url.QueryEscape(medTitle)

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", "https://eservice.gu.spb.ru/portalFront/proxy/async?filter="+hh+"&operation=getMedicament", nil,
	)
	// добавляем заголовки
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:74.0) Gecko/20100101 Firefox/74.0")

	resp, err := client.Do(req)
	if err != nil {
		return Jsn{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("код ответа: %v", resp.StatusCode)
		return Jsn{}, err
	}

	j := Jsn{}
	data, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(data, &j)

	return j, nil
}

// IsErrExistInJSON - Проверяет json структуру на наличие ошибок
func IsErrExistInJSON(j Jsn) bool {
	if u := j.Errors; u != "" {
		return true
	}
	return false
}

// ParseJSON - парсит json структуру и возвращает готовый текст сообщения для пользователя
func ParseJSON(j Jsn) string {
	var text []string

	title := fmt.Sprintln("Название: ", j.Model.Result[0].Name)

	for name := range j.Model.Result[0].Districts {
		district := fmt.Sprint("\n\n[[", j.Model.Result[0].Districts[name].Name, " ]]\n\n")
		text = append(text, district)

		for _, apothecary := range j.Model.Result[0].Districts[name].Apothecaries {
			name := apothecary.Name
			addr := apothecary.Addr
			a := strings.Trim(addr, "  * На момент обращения в аптеку не гарантируется наличие лекарственного препарата к выдаче, в связи с ограничением количества препарата в аптеке. Информацию о наличии препарата необходимо уточнить по телефону")

			s := strings.Split(a, ",")

			// index := fmt.Sprintln("Индекс: ", s[0])
			street := strings.TrimPrefix(s[2], " ")
			house := s[3]
			address := street + " " + house

			fedExemption := fmt.Sprintln("Федеральная льгота: ", apothecary.FedExemption)
			//Региональная льгота
			regExemption := fmt.Sprintln("Региональнальная льгота: ", apothecary.RegExemption)
			//Писхиатрическая льгота
			psyExemption := fmt.Sprintln("Психиатрическая льгота: ", apothecary.PsyExemption)
			//ВЗН
			vzn := fmt.Sprintln("ВЗН: ", apothecary.VZN)

			apoth := name + "\n" + address + "\n\n" + fedExemption + regExemption + psyExemption + vzn

			text = append(text, apoth)
		}
	}

	msg := title
	for _, i := range text {
		msg += i
	}

	return msg
}

// CyclicMedSearch - Проверяет список подписок, после чего опрашивает Гос Услуги
// на наличие эти лекарств в городе. Полученный результат сравнивается с
// состоянием в базе данных. если значение Avaliability сменяется на true,
// то пользователю отправляется уведомление.
func CyclicMedSearch(bot *tgbotapi.BotAPI, c chan bool) {

	conn, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// Проверяем наличе хотя бы одной подписки, дабы избежать ошибок, связанных
	// с попыткой чтения не существующей информации
	anySub, err := db.AreTheAnySubscriptions(conn)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if anySub == true {

		medsID, err := db.FindAllSubscriptionsMed(conn)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		for _, id := range medsID {
			title, err := db.FindMedTitle(conn, id)
			if err != nil {
				log.Fatalf("%+v", err)
			}

			medResp, err := ReqMedInfo(title)
			if err != nil {
				log.Fatalf("%+v", err)
			}

			isErr := IsErrExistInJSON(medResp)
			if err != nil {
				log.Fatalf("%+v", err)
			}

			availabillity, err := db.CheckAvailability(conn, id)
			if err != nil {
				log.Fatalf("%+v", err)
			}

			// Наличие ошибки говорит нам о том, что лекарства в данный момент нигде нет
			// Значит, мы проверяем значение Availability в базе данных
			if isErr == true && availabillity == true {
				db.ChangeAvailability(conn, id, false)
				// Ставим заметку о дате, когда лекарство закончилось в городе
			}

			if isErr == false && availabillity == false {
				db.ChangeAvailability(conn, id, true)
				// Теперь нужно уведомить всех пользователей,
				// которые подписаны на данное лекарство
				users, err := db.FindWhoSubToMed(conn, id)
				if err != nil {
					log.Fatalf("%+v", err)
				}

				for _, user := range users {
					chatID, err := db.FindChatID(conn, user)
					if err != nil {
						log.Fatalf("%+v", err)
					}

					msgText := ParseJSON(medResp)

					msgConf := tgbotapi.NewMessage(int64(chatID), msgText)
					msgConf.ReplyMarkup = keyboards.ViewMedKeyboard

					if _, err := bot.Send(msgConf); err != nil {
						log.Fatalf("%+v", err)
					}
				}
			}
		}
	}
}

// ReadFileWithMeds - считывает данные из файла drugs.txt и подготавливает их для
// передачи в функцию InitMedList, котрая заполнит ими базу данных
func ReadFileWithMeds() ([]string, error) {
	file, err := os.Open("drugs.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
