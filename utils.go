package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func initEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	botToken = os.Getenv("BOT_TOKEN")
	baseUrl = os.Getenv("BASE_URL")

	return nil
}

func getUpdates(apiUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(apiUrl + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		fmt.Println("Something went wrong in requset: ", err)
		return nil, err
	}

	// ответ от сервера получаем в байтах, необходимо обработать его
	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		if err != nil {
			fmt.Println("Something went wrong in error handling: ", err)
			return nil, err
		}
	}

	var restResponse RestResponse

	// необходим распарсить json, который получили от сервера, который приведем к структуре RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Println("Something went wrong in parse json: ", err)
		return nil, err
	}

	return restResponse.Result, nil
}

func checkMainCommand(message Message) {
	if message.Text == "/start" {
		fmt.Println("Hi, nice to meet you!")
	} else if message.Text == "" {
		fmt.Println("Bye bye! see you soon")
	} else {
		fmt.Println("Unknow command =/")
	}
}
