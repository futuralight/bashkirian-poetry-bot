package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

const (
	//PoemistAPIUrl - Poemist api url
	PoemistAPIUrl = "https://www.poemist.com/api/v1/randompoems"
)

//PoemistItem - JSON item from poemist api
type PoemistItem struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Poet    struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
}

var Keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/poem"),
	),
)

type YandexTranslateResponse struct {
	Code int      `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}

func init() {
	loadEnv()
}

func loadEnv() error {
	godotenv.Load()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	err = godotenv.Load(dir + "/.env") //Загрузка .env файла
	return nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELE_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = "Здарова"
			case "poem":
				msg.ReplyMarkup = Keyboard
				msg.Text = getBashPoem()
			default:
				msg.Text = "Пиши /poem"
			}
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пиши /poem")
			bot.Send(msg)
		}
	}
}

func getBashPoem() string {
	poems, err := getPoems()
	if err != nil {
		panic(err)
	}
	text := translateRequest(poems[0].Content)
	return text
}

func translateRequest(text string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://translate.yandex.net/api/v1.5/tr.json/translate", nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	q.Add("key", os.Getenv("YANDEX_TRANSLATE_TOKEN"))
	q.Add("lang", "en-ba")
	q.Add("format", "plain")
	q.Add("text", text)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	yaResponse := YandexTranslateResponse{}
	err = json.Unmarshal([]byte(rawData), &yaResponse)
	if err != nil {
		panic(err)

	}
	return yaResponse.Text[0]
}

func getPoems() ([]PoemistItem, error) {
	resp, err := http.Get(PoemistAPIUrl)
	if err != nil {
		return nil, err
	}
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	poems := []PoemistItem{}
	err = json.Unmarshal([]byte(rawData), &poems)
	if err != nil {
		return nil, err
	}
	return poems, nil
}
