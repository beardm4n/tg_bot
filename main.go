package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	BASE_URL = "https://api.telegram.org/bot"
)

var (
	telegramMethods = map[string]string{
		"GET_UPDATES":  "getUpdates",
		"SEND_MESSAGE": "sendMessage",
	}
	offset   = 0
	botToken string
)

func main() {
	err := initEnv()
	if err != nil {
		log.Fatal("Error in loadEnv method .env: ", err)
	}

	apiUrl := BASE_URL + botToken + "/" + telegramMethods["GET_UPDATES"]

	for {
		updates, err := getUpdates(apiUrl, offset)
		if err != nil {
			log.Println("Something went wrong in getUpdates", err.Error())
		}

		for _, update := range updates {
			offset = update.UpdateId + 1
		}

		fmt.Printf("Updates: %v\n", updates)

		// Пауза между запросами на получение обновлений
		time.Sleep(time.Second)
	}
}

func getUpdates(apiUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(apiUrl + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		fmt.Println("Something went wrong in requset:", err)
		return nil, err
	}

	// ответ от сервера получаем в байтах, необходимо обработать его
	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		if err != nil {
			fmt.Println("Something went wrong in error handling:", err)
			return nil, err
		}
	}

	var restResponse RestResponse

	// необходим распарсить json, который получили от сервера, который приведем к структуре RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Println("Something went wrong in parse json:", err)
		return nil, err
	}
	
	return restResponse.Result, nil
}

func initEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	botToken = os.Getenv("BOT_TOKEN")

	return nil
}
