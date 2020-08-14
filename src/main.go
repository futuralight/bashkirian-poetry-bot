package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

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
	poems, err := getPoems()
	if err != nil {
		panic(err)
	}
	// fmt.Println(poems[0].Content)
	translateRequest(poems[0].Content)
}

func translateRequest(text string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://translate.yandex.net/api/v1.5/tr.json/translate", nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	q.Add("key", os.Getenv("YANDEX_TRANSLATE_TOKEN"))
	q.Add("lang", "ru-ba")
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
	fmt.Println(req.URL.String())
	// fmt.Println(yaResponse)
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
